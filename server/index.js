import express from 'express';
import cors from 'cors';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs-extra';
import { v4 as uuidv4 } from 'uuid';
import { createWriteStream } from 'fs';
import { pipeline } from 'stream/promises';
import SourceCollector from './sourcecollector.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = express();
const PORT = process.env.PORT || 3001;

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.static(join(__dirname, '../dist')));

// Storage directories
const TEMP_DIR = join(__dirname, '../temp');
const OUTPUT_DIR = join(__dirname, '../output');

// Ensure directories exist
await fs.ensureDir(TEMP_DIR);
await fs.ensureDir(OUTPUT_DIR);

// Store job status
const jobs = new Map();

// Helper function to validate GitHub URL
function validateGitHubUrl(url) {
  const githubRegex = /^https:\/\/github\.com\/[a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+(?:\.git)?$/;
  return githubRegex.test(url);
}

// Helper function to extract repo info from URL
function extractRepoInfo(url) {
  const match = url.match(/github\.com\/([^\/]+)\/([^\/]+?)(?:\.git)?$/);
  if (match) {
    return {
      owner: match[1],
      repo: match[2]
    };
  }
  return null;
}

// Helper function to download and extract GitHub repository
async function downloadRepository(repoInfo, extractPath) {
  const { owner, repo } = repoInfo;
  
  // Try main branch first, then master
  const branches = ['main', 'master'];
  let lastError;
  
  for (const branch of branches) {
    try {
      const archiveUrl = `https://github.com/${owner}/${repo}/archive/refs/heads/${branch}.zip`;
      const response = await fetch(archiveUrl);
      
      if (response.ok) {
        await extractZipFromResponse(response, extractPath, `${repo}-${branch}`);
        return extractPath;
      }
    } catch (error) {
      lastError = error;
      continue;
    }
  }
  
  throw new Error(`Repository not found or not accessible. Last error: ${lastError?.message || 'Unknown error'}`);
}

// Helper function to extract zip from response
async function extractZipFromResponse(response, extractPath, expectedFolderName) {
  const zipPath = `${extractPath}.zip`;
  
  try {
    // Download zip file
    const fileStream = createWriteStream(zipPath);
    await pipeline(response.body, fileStream);
    
    // Try to use adm-zip if available
    let AdmZip;
    try {
      AdmZip = (await import('adm-zip')).default;
    } catch (error) {
      AdmZip = null;
    }
    
    if (AdmZip) {
      // Use adm-zip
      const zip = new AdmZip(zipPath);
      zip.extractAllTo(dirname(extractPath), true);
      
      // Find the extracted folder and rename it to our expected path
      const extractedFolder = join(dirname(extractPath), expectedFolderName);
      if (await fs.pathExists(extractedFolder)) {
        await fs.move(extractedFolder, extractPath);
      }
    } else {
      // Fallback: manual extraction using Node.js streams and basic zip parsing
      throw new Error('ZIP extraction not available. Please install adm-zip dependency.');
    }
  } finally {
    // Clean up zip file
    try {
      await fs.unlink(zipPath);
    } catch (error) {
      // Ignore cleanup errors
    }
  }
}

// API Routes

// Get all jobs
app.get('/api/jobs', (req, res) => {
  const jobList = Array.from(jobs.entries()).map(([id, job]) => ({
    id,
    ...job
  }));
  res.json(jobList);
});

// Get specific job
app.get('/api/jobs/:id', (req, res) => {
  const job = jobs.get(req.params.id);
  if (!job) {
    return res.status(404).json({ error: 'Job not found' });
  }
  res.json({ id: req.params.id, ...job });
});

// Clone repository and process
app.post('/api/clone', async (req, res) => {
  const { repoUrl } = req.body;

  if (!repoUrl) {
    return res.status(400).json({ error: 'Repository URL is required' });
  }

  if (!validateGitHubUrl(repoUrl)) {
    return res.status(400).json({ error: 'Invalid GitHub URL' });
  }

  const jobId = uuidv4();
  const repoInfo = extractRepoInfo(repoUrl);
  
  if (!repoInfo) {
    return res.status(400).json({ error: 'Could not parse repository information' });
  }

  const job = {
    status: 'pending',
    repoUrl,
    repoInfo,
    createdAt: new Date().toISOString(),
    progress: 0,
    message: 'Initializing...'
  };

  jobs.set(jobId, job);

  // Start processing in background
  processRepository(jobId, repoUrl, repoInfo);

  res.json({ jobId, status: 'started' });
});

// Download processed file
app.get('/api/download/:id', async (req, res) => {
  const job = jobs.get(req.params.id);
  
  if (!job) {
    return res.status(404).json({ error: 'Job not found' });
  }

  if (job.status !== 'completed') {
    return res.status(400).json({ error: 'Job not completed yet' });
  }

  const filePath = job.outputPath;
  
  if (!fs.existsSync(filePath)) {
    return res.status(404).json({ error: 'Output file not found' });
  }

  const filename = `${job.repoInfo.owner}-${job.repoInfo.repo}-source.txt`;
  
  res.download(filePath, filename, (err) => {
    if (err) {
      console.error('Download error:', err);
      res.status(500).json({ error: 'Download failed' });
    }
  });
});

// Delete job and cleanup
app.delete('/api/jobs/:id', async (req, res) => {
  const job = jobs.get(req.params.id);
  
  if (!job) {
    return res.status(404).json({ error: 'Job not found' });
  }

  // Cleanup files
  try {
    if (job.clonePath && fs.existsSync(job.clonePath)) {
      await fs.remove(job.clonePath);
    }
    if (job.outputPath && fs.existsSync(job.outputPath)) {
      await fs.remove(job.outputPath);
    }
  } catch (error) {
    console.error('Cleanup error:', error);
  }

  jobs.delete(req.params.id);
  res.json({ success: true });
});

// Background processing function
async function processRepository(jobId, repoUrl, repoInfo) {
  const job = jobs.get(jobId);
  
  try {
    // Update status
    job.status = 'cloning';
    job.progress = 10;
    job.message = 'Downloading repository...';

    // Create unique directory for this job
    const clonePath = join(TEMP_DIR, `${jobId}-${repoInfo.repo}`);
    const outputPath = join(OUTPUT_DIR, `${jobId}-source.txt`);

    job.clonePath = clonePath;
    job.outputPath = outputPath;

    // Download repository using GitHub API
    await downloadRepository(repoInfo, clonePath);

    job.progress = 50;
    job.message = 'Repository downloaded, processing files...';

    // Run JavaScript-based SourceCollector
    const collector = new SourceCollector(clonePath, outputPath);
    const result = await collector.run();

    // Get file stats
    const stats = await fs.stat(outputPath);
    
    job.status = 'completed';
    job.progress = 100;
    job.message = `Processing completed successfully (${result.filesProcessed} files processed)`;
    job.completedAt = new Date().toISOString();
    job.fileSize = stats.size;
    job.filesProcessed = result.filesProcessed;

  } catch (error) {
    console.error(`Job ${jobId} failed:`, error);
    job.status = 'failed';
    job.error = error.message;
    job.message = `Failed: ${error.message}`;
  }
}

// Serve React app for all other routes
app.get('*', (req, res) => {
  res.sendFile(join(__dirname, '../dist/index.html'));
});

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
  console.log('Using JavaScript-based SourceCollector (Go binary not required)');
});