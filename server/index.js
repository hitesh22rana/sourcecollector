import express from 'express';
import cors from 'cors';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs-extra';
import { v4 as uuidv4 } from 'uuid';
import { spawn } from 'child_process';
import archiver from 'archiver';
import { createWriteStream } from 'fs';
import { pipeline } from 'stream/promises';

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
  const archiveUrl = `https://github.com/${owner}/${repo}/archive/refs/heads/main.zip`;
  
  try {
    // Try main branch first
    const response = await fetch(archiveUrl);
    if (!response.ok) {
      // If main branch doesn't exist, try master branch
      const masterUrl = `https://github.com/${owner}/${repo}/archive/refs/heads/master.zip`;
      const masterResponse = await fetch(masterUrl);
      if (!masterResponse.ok) {
        throw new Error(`Repository not found or not accessible: ${response.status}`);
      }
      return await extractZipFromResponse(masterResponse, extractPath, `${repo}-master`);
    }
    return await extractZipFromResponse(response, extractPath, `${repo}-main`);
  } catch (error) {
    throw new Error(`Failed to download repository: ${error.message}`);
  }
}

// Helper function to extract zip from response
async function extractZipFromResponse(response, extractPath, expectedFolderName) {
  const zipPath = `${extractPath}.zip`;
  
  // Download zip file
  const fileStream = createWriteStream(zipPath);
  await pipeline(response.body, fileStream);
  
  // Extract zip file using Node.js built-in modules
  let AdmZip;
  try {
    AdmZip = (await import('adm-zip')).default;
  } catch (error) {
    AdmZip = null;
  }
  
  if (!AdmZip) {
    // Fallback: use unzip command if available
    return new Promise((resolve, reject) => {
      const unzipProcess = spawn('unzip', ['-q', zipPath, '-d', dirname(extractPath)]);
      
      unzipProcess.on('close', (code) => {
        // Clean up zip file
        fs.unlink(zipPath).catch(() => {});
        
        if (code === 0) {
          // Find the extracted folder and rename it
          const extractedFolder = join(dirname(extractPath), expectedFolderName);
          if (fs.existsSync(extractedFolder)) {
            fs.rename(extractedFolder, extractPath).then(resolve).catch(reject);
          } else {
            resolve(extractPath);
          }
        } else {
          reject(new Error(`Unzip failed with code ${code}`));
        }
      });
      
      unzipProcess.on('error', reject);
    });
  }
  
  // Use adm-zip if available
  try {
    const zip = new AdmZip(zipPath);
    zip.extractAllTo(dirname(extractPath), true);
    
    // Clean up zip file
    await fs.unlink(zipPath);
    
    // Find the extracted folder and rename it to our expected path
    const extractedFolder = join(dirname(extractPath), expectedFolderName);
    if (await fs.pathExists(extractedFolder)) {
      await fs.move(extractedFolder, extractPath);
    }
    
    return extractPath;
  } catch (error) {
    // Clean up zip file on error
    await fs.unlink(zipPath).catch(() => {});
    throw error;
  }
}

// Helper function to ensure SourceCollector binary exists and is executable
async function ensureSourceCollectorBinary() {
  const sourcecollectorPath = join(__dirname, '../bin/sourcecollector');
  
  // Check if binary exists
  if (!await fs.pathExists(sourcecollectorPath)) {
    throw new Error('SourceCollector binary not found. Please build the Go application first by running: go build -o bin/sourcecollector .');
  }
  
  // Make sure it's executable
  try {
    await fs.chmod(sourcecollectorPath, 0o755);
  } catch (error) {
    console.warn('Could not set executable permissions:', error.message);
  }
  
  return sourcecollectorPath;
}

// Helper function to run sourcecollector
async function runSourceCollector(inputPath, outputPath) {
  const sourcecollectorPath = await ensureSourceCollectorBinary();
  
  return new Promise((resolve, reject) => {
    const process = spawn(sourcecollectorPath, [
      '--input', inputPath,
      '--output', outputPath,
      '--fast'
    ], {
      stdio: ['pipe', 'pipe', 'pipe']
    });

    let stdout = '';
    let stderr = '';

    process.stdout.on('data', (data) => {
      stdout += data.toString();
    });

    process.stderr.on('data', (data) => {
      stderr += data.toString();
    });

    process.on('close', (code) => {
      if (code === 0) {
        resolve({ stdout, stderr });
      } else {
        reject(new Error(`SourceCollector failed with code ${code}: ${stderr || stdout}`));
      }
    });

    process.on('error', (error) => {
      reject(new Error(`Failed to execute SourceCollector: ${error.message}`));
    });
  });
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

    // Run sourcecollector
    await runSourceCollector(clonePath, outputPath);

    // Get file stats
    const stats = await fs.stat(outputPath);
    
    job.status = 'completed';
    job.progress = 100;
    job.message = 'Processing completed successfully';
    job.completedAt = new Date().toISOString();
    job.fileSize = stats.size;

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
});