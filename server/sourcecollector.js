import fs from 'fs-extra';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// File extensions that should be processed
const VALID_EXTENSIONS = new Set([
  '.js', '.jsx', '.ts', '.tsx', '.vue', '.svelte',
  '.py', '.rb', '.php', '.java', '.c', '.cpp', '.h', '.hpp',
  '.cs', '.go', '.rs', '.swift', '.kt', '.scala',
  '.html', '.htm', '.css', '.scss', '.sass', '.less',
  '.json', '.xml', '.yaml', '.yml', '.toml',
  '.sql', '.sh', '.bash', '.zsh', '.fish',
  '.r', '.m', '.pl', '.lua', '.dart', '.elm',
  '.clj', '.cljs', '.hs', '.ml', '.fs', '.ex', '.exs',
  '.jl', '.nim', '.cr', '.zig', '.odin', '.v',
  '.dockerfile', '.makefile', '.cmake', '.gradle',
  '.config', '.conf', '.ini', '.env'
]);

// Files and directories to ignore
const IGNORE_PATTERNS = [
  'node_modules', 'bower_components', 'dist', 'build', '.next', '.nuxt',
  'coverage', 'vendor', 'tmp', 'temp', 'logs', 'log', '.git', '.svn',
  '.DS_Store', 'Thumbs.db', 'package-lock.json', 'yarn.lock',
  '.vscode', '.idea', '__pycache__', '.pytest_cache', '.mypy_cache',
  'target', 'bin', 'obj', '.gradle', '.maven'
];

// Binary file extensions to ignore
const BINARY_EXTENSIONS = new Set([
  '.exe', '.dll', '.so', '.dylib', '.a', '.lib', '.o', '.obj',
  '.jpg', '.jpeg', '.png', '.gif', '.bmp', '.svg', '.ico',
  '.mp3', '.mp4', '.avi', '.mov', '.wav', '.flac',
  '.pdf', '.doc', '.docx', '.xls', '.xlsx', '.ppt', '.pptx',
  '.zip', '.tar', '.gz', '.rar', '.7z', '.dmg', '.iso'
]);

class SourceCollector {
  constructor(inputPath, outputPath) {
    this.inputPath = inputPath;
    this.outputPath = outputPath;
    this.files = [];
  }

  shouldIgnore(filePath) {
    const relativePath = path.relative(this.inputPath, filePath);
    const parts = relativePath.split(path.sep);
    
    // Check if any part of the path matches ignore patterns
    for (const part of parts) {
      if (IGNORE_PATTERNS.some(pattern => part.includes(pattern))) {
        return true;
      }
      // Ignore hidden files and directories (starting with .)
      if (part.startsWith('.') && part !== '.' && part !== '..') {
        return true;
      }
    }
    
    return false;
  }

  isValidFile(filePath) {
    const ext = path.extname(filePath).toLowerCase();
    
    // Ignore binary files
    if (BINARY_EXTENSIONS.has(ext)) {
      return false;
    }
    
    // Check if it's a valid programming file
    if (VALID_EXTENSIONS.has(ext)) {
      return true;
    }
    
    // Check for files without extensions that might be important
    const basename = path.basename(filePath).toLowerCase();
    const importantFiles = [
      'dockerfile', 'makefile', 'rakefile', 'gemfile', 'procfile',
      'readme', 'license', 'changelog', 'contributing', 'authors'
    ];
    
    return importantFiles.some(name => basename.startsWith(name));
  }

  async scanDirectory(dirPath) {
    try {
      const entries = await fs.readdir(dirPath, { withFileTypes: true });
      
      for (const entry of entries) {
        const fullPath = path.join(dirPath, entry.name);
        
        if (this.shouldIgnore(fullPath)) {
          continue;
        }
        
        if (entry.isDirectory()) {
          await this.scanDirectory(fullPath);
        } else if (entry.isFile() && this.isValidFile(fullPath)) {
          this.files.push(fullPath);
        }
      }
    } catch (error) {
      console.warn(`Warning: Could not read directory ${dirPath}:`, error.message);
    }
  }

  generateTree() {
    const tree = {};
    
    for (const filePath of this.files) {
      const relativePath = path.relative(this.inputPath, filePath);
      const parts = relativePath.split(path.sep);
      
      let current = tree;
      for (let i = 0; i < parts.length - 1; i++) {
        const part = parts[i];
        if (!current[part]) {
          current[part] = {};
        }
        current = current[part];
      }
      
      const fileName = parts[parts.length - 1];
      current[fileName] = filePath;
    }
    
    return tree;
  }

  formatTree(tree, prefix = '', isLast = true) {
    let result = '';
    const entries = Object.entries(tree);
    
    entries.forEach(([name, value], index) => {
      const isLastEntry = index === entries.length - 1;
      const connector = isLastEntry ? '└── ' : '├── ';
      result += prefix + connector + name + '\n';
      
      if (typeof value === 'object' && value !== null) {
        const newPrefix = prefix + (isLastEntry ? '    ' : '│   ');
        result += this.formatTree(value, newPrefix, isLastEntry);
      }
    });
    
    return result;
  }

  async processFiles() {
    let content = '';
    
    // Generate and add file tree
    const tree = this.generateTree();
    content += 'Source code files structure\n\n';
    content += this.formatTree(tree);
    content += '\n\n';
    
    // Process each file
    for (const filePath of this.files) {
      try {
        const relativePath = path.relative(this.inputPath, filePath);
        const fileName = path.basename(filePath);
        const fileContent = await fs.readFile(filePath, 'utf8');
        
        content += `Name: ${fileName}\n`;
        content += `Path: ${relativePath}\n`;
        content += '```\n';
        content += fileContent;
        content += '\n```\n\n';
      } catch (error) {
        console.warn(`Warning: Could not read file ${filePath}:`, error.message);
        const relativePath = path.relative(this.inputPath, filePath);
        const fileName = path.basename(filePath);
        
        content += `Name: ${fileName}\n`;
        content += `Path: ${relativePath}\n`;
        content += '```\n';
        content += `[Error reading file: ${error.message}]\n`;
        content += '```\n\n';
      }
    }
    
    return content;
  }

  async run() {
    console.log('Starting SourceCollector...');
    
    // Check if input path exists
    if (!await fs.pathExists(this.inputPath)) {
      throw new Error(`Input path does not exist: ${this.inputPath}`);
    }
    
    // Scan directory for files
    console.log('Scanning directory...');
    await this.scanDirectory(this.inputPath);
    
    console.log(`Found ${this.files.length} files to process`);
    
    // Process files and generate content
    console.log('Processing files...');
    const content = await this.processFiles();
    
    // Write output file
    console.log('Writing output file...');
    await fs.ensureDir(path.dirname(this.outputPath));
    await fs.writeFile(this.outputPath, content, 'utf8');
    
    console.log(`SourceCollector completed successfully!`);
    console.log(`Output written to: ${this.outputPath}`);
    
    return {
      filesProcessed: this.files.length,
      outputSize: content.length
    };
  }
}

export default SourceCollector;