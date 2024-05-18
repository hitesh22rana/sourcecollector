package pkg

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	unwantedFilesAndFolders = []string{
		// General
		"node_modules",     // Node.js dependencies
		"bower_components", // Bower dependencies
		"dist",             // Distribution directory
		"build",            // Build directory
		".next",            // Next.js build directory
		".nuxt",            // Nuxt.js build directory
		"databases",        // Database files
		"data",             // Data files
		"logs",             // Log files
		"out",              // Next.js build directory
		"public",           // Public directory
		"coverage",         // Code coverage reports
		"vendor",           // Composer dependencies
		"tmp",              // Temporary files
		"temp",             // Temporary files
		".log",             // Log files
		".tmp",             // Temporary files
		".bak",             // Backup files
		".swp",             // Swap files
		"LCK..",            // Lock files
		".DS_Store",        // macOS file system metadata
		"Thumbs.db",        // Windows file system metadata
		"LICENSE",          // License files
		"AUTHORS",          // Authors files
		"CONTRIBUTORS",     // Contributors files
		"CHANGELOG",        // Changelog files
		"CHANGES",          // Changes files
		"HISTORY",          // History files
		"NOTICE",           // Notice files
		"README",           // Readme files
		"TODO",             // Todo files

		// JavaScript
		"package-lock.json", // NPM lock file
		"yarn.lock",         // Yarn lock file
		"jest.config.js",    // Jest configuration
		"jest.setup.js",     // Jest setup
		"jest.json",         // Jest configuration
		"jest",              // Jest configuration
		"webpack.config.js", // Webpack configuration
		"rollup.config.js",  // Rollup configuration
		"gulpfile.js",       // Gulp configuration
		"Gruntfile.js",      // Grunt configuration
		"tsconfig.json",     // TypeScript configuration
		"tslint.json",       // TypeScript lint configuration
		"jsconfig.json",     // JavaScript configuration
		"babel.config.js",   // Babel configuration
		"prettier.config",   // Prettier configuration

		// C/C++
		".o",   // Object files
		".obj", // Object files
		".so",  // Shared objects
		".a",   // Static libraries
		".lib", // Static libraries
		".dll", // Dynamic Link Libraries
		".exe", // Executable files
		".out", // Executable files

		// Java
		".class", // Compiled Java classes
		".jar",   // Java Archive
		".war",   // Web Application Archive
		".ear",   // Enterprise Application Archive

		// Python
		".pyc",          // Compiled Python files
		".pyo",          // Optimized Python files
		"__pycache__",   // Python cache directory
		".egg-info",     // Python package metadata
		".eggs",         // Python package
		".pytest_cache", // Pytest cache directory
		".tox",          // Tox directory
		".mypy_cache",   // Mypy cache directory
		".hypothesis",   // Hypothesis directory
		".nox",          // Nox directory
		".coverage",     // Coverage report
		".cache",        // Cache directory
		".env",          // Virtual environment
		".venv",         // Virtual environment
		"/venv",         // Virtual environment
		"venv",          // Virtual environment

		// Ruby
		".gem",    // RubyGem package
		".bundle", // Bundler directory

		// Go
		".exe",   // Executable files
		".test",  // Test binary
		".out",   // Output files
		"go.sum", // Go dependencies
		"go.mod", // Go dependencies

		// Rust
		"target", // Cargo build directory

		// Swift
		".xcodeproj",   // Xcode project
		".xcworkspace", // Xcode workspace
		".swiftmodule", // Compiled Swift module
		".swiftdoc",    // Swift documentation

		// .NET
		"bin", // Binary output directory
		"obj", // Object output directory

		// Version Control
		".git",           // Git directory
		".gitignore",     // Git ignore file
		".gitattributes", // Git attributes file
		".gitmodules",    // Git modules file
		".svn",           // Subversion directory
		".hg",            // Mercurial directory

		// Editor and IDE specific
		".DS_Store",     // macOS file system metadata
		"Thumbs.db",     // Windows file system metadata
		".idea",         // IntelliJ IDEA project files
		".iml",          // IntelliJ IDEA module files
		".vscode",       // Visual Studio Code settings
		".suo",          // Visual Studio solution user options
		".user",         // Visual Studio user options
		".userosscache", // Visual Studio user options
		".sln",          // Visual Studio solution files
		".psess",        // Visual Studio performance session files
		".vsp",          // Visual Studio performance report files

		// Miscellaneous
		".iso", // Disk images
		".tar", // Tarballs
		".gz",  // Gzip compressed files
		".zip", // Zip compressed files
		".7z",  // 7-Zip compressed files
		".rar", // RAR compressed files
	}

	validProgrammingFileExtensions = []string{
		".1c",          // 1C Enterprise
		".4th",         // Forth
		".6pl",         // Perl
		".6pm",         // Perl
		".aba",         // Ada
		".adb",         // Ada
		".ads",         // Ada
		".agc",         // AGC
		".ahk",         // AutoHotkey
		".aj",          // AspectJ
		".als",         // Alloy
		".apl",         // APL
		".applescript", // AppleScript
		".arc",         // Arc
		".as",          // ActionScript
		".asm",         // Assembly
		".asp",         // ASP
		".aspx",        // ASP.NET
		".awk",         // AWK
		".bas",         // BASIC
		".bash",        // Bash
		".bat",         // Batch file
		".bb",          // BlitzBasic
		".bbx",         // Berry
		".bdf",         // BDF font
		".bf",          // Brainfuck
		".bmx",         // BlitzMax
		".boo",         // Boo
		".brs",         // BrightScript
		".bsv",         // Bluespec SystemVerilog
		".c",           // C
		".c++",         // C++
		".cbl",         // COBOL
		".cc",          // C++
		".ceylon",      // Ceylon
		".chpl",        // Chapel
		".cjs",         // CommonJS
		".cl",          // Common Lisp
		".clj",         // Clojure
		".cljs",        // ClojureScript
		".cls",         // Visual Basic
		".cmake",       // CMake
		".cob",         // COBOL
		".coffee",      // CoffeeScript
		".cp",          // C++
		".cpp",         // C++
		".cpy",         // Python
		".cr",          // Crystal
		".cs",          // C#
		".csh",         // C Shell
		".cson",        // CoffeeScript Object Notation
		".csproj",      // C#
		".css",         // CSS
		".cu",          // CUDA
		".cxx",         // C++
		".d",           // D
		".dart",        // Dart
		".dats",        // ATS
		".dbs",         // SQL
		".dcl",         // Clean
		".decls",       // Clean
		".diderot",     // Diderot
		".dita",        // DITA
		".ditamap",     // DITA
		".djt",         // D
		".dml",         // D
		".doh",         // D
		".dot",         // Graphviz
		".dpr",         // Delphi
		".druby",       // dRuby
		".dtx",         // LaTeX
		".dylan",       // Dylan
		".dyl",         // Dylan
		".e",           // Eiffel
		".ec",          // C
		".eh",          // C
		".el",          // Emacs Lisp
		".elm",         // Elm
		".em",          // E
		".erl",         // Erlang
		".ex",          // Elixir
		".exs",         // Elixir
		".f",           // Fortran
		".f90",         // Fortran
		".f95",         // Fortran
		".factor",      // Factor
		".fan",         // Fantom
		".fth",         // Forth
		".fish",        // fish shell
		".for",         // Fortran
		".forth",       // Forth
		".fs",          // F#
		".fsi",         // F#
		".fsscript",    // F#
		".fsx",         // F#
		".g",           // G-code
		".gap",         // GAP
		".gawk",        // AWK
		".gdb",         // GDB
		".gd",          // GDScript
		".gdns",        // Godot
		".ged",         // Godot
		".glf",         // GLSL
		".gml",         // GameMaker Language
		".go",          // Go
		".gs",          // Google Apps Script
		".gsp",         // Groovy Server Pages
		".gst",         // GAMS
		".gsx",         // GAMS
		".gvy",         // Groovy
		".h",           // C/C++
		".hack",        // Hack
		".haml",        // Haml
		".handlebars",  // Handlebars
		".hbs",         // Handlebars
		".hs",          // Haskell
		".html",        // HTML
		".htm",         // HTML
		".hx",          // Haxe
		".hxx",         // C++
		".ice",         // ICE
		".iced",        // IcedCoffeeScript
		".idr",         // Idris
		".ijs",         // J
		".imba",        // Imba
		".inc",         // PHP
		".ini",         // Configuration file
		".ino",         // Arduino
		".io",          // Io
		".j",           // Java
		".jade",        // Jade
		".java",        // Java
		".jl",          // Julia
		".js",          // JavaScript
		".jsb",         // JavaScript
		".jscad",       // OpenJSCAD
		".jsfl",        // JavaScript
		".jsh",         // JavaScript
		".json",        // JSON
		".json5",       // JSON5
		".jsx",         // JavaScript
		".jflex",       // JFlex
		".jison",       // Jison
		".jisonlex",    // Jison Lex
		".jl",          // Julia
		".kak",         // Kakoune
		".kicad_pcb",   // KiCad
		".kicad_sch",   // KiCad
		".kit",         // Kite
		".kt",          // Kotlin
		".kts",         // Kotlin
		".kxi",         // Kite
		".kxml",        // Kite
		".l",           // Lisp
		".lagda",       // Agda
		".lagda.md",    // Agda
		".lagda.rst",   // Agda
		".lean",        // Lean
		".less",        // LESS
		".lhs",         // Literate Haskell
		".lid",         // D
		".lisp",        // Lisp
		".lkt",         // Inkling
		".lmo",         // Limbo
		".lua",         // Lua
		".ly",          // LilyPond
		".m",           // Objective-C
		".mac",         // M4
		".mak",         // Makefile
		".make",        // Makefile
		".man",         // Unix Manual
		".markdown",    // Markdown
		".marko",       // Marko
		".mat",         // MATLAB
		".mata",        // MATLAB
		".matlab",      // MATLAB
		".maxpat",      // Max
		".md",          // Markdown
		".mediawiki",   // MediaWiki
		".mirah",       // Mirah
		".mjml",        // MJML
		".mjs",         // JavaScript
		".ml",          // OCaml
		".mli",         // OCaml
		".mo",          // Modula-2
		".monkey",      // Monkey
		".moon",        // MoonScript
		".ms",          // Common Lisp
		".mumps",       // MUMPS
		".mustache",    // Mustache
		".mxml",        // Flex
		".n",           // N
		".nawk",        // AWK
		".nb",          // Mathematica
		".ncl",         // NCL
		".nl",          // GAMS
		".nix",         // Nix
		".numpy",       // Python
		".nu",          // Nu
		".num",         // Python
		".nut",         // Squirrel
		".o",           // Object file
		".obj",         // Object file
		".odin",        // Odin
		".omgrofl",     // Omgrofl
		".org",         // Org-mode
		".ox",          // Ox
		".p",           // Pascal
		".p6",          // Perl 6
		".pac",         // JavaScript
		".parrot",      // Parrot
		".pas",         // Pascal
		".patch",       // Patch
		".pat",         // D
		".pawn",        // Pawn
		".pbf",         // D
		".pbi",         // PureBasic
		".pde",         // Processing
		".perl",        // Perl
		".php",         // PHP
		".phps",        // PHP
		".phtml",       // PHP
		".pig",         // Pig
		".pike",        // Pike
		".pl",          // Perl
		".pl6",         // Perl 6
		".pls",         // PLSQL
		".plx",         // Perl
		".pm",          // Perl
		".pml",         // PHP
		".pm6",         // Perl 6
		".pmod",        // D
		".pod",         // Perl
		".pony",        // Pony
		".pp",          // Puppet
		".prg",         // FoxPro
		".pro",         // Prolog
		".prolog",      // Prolog
		".ps1",         // PowerShell
		".psc1",        // PowerShell
		".psm1",        // PowerShell
		".purs",        // PureScript
		".py",          // Python
		".py3",         // Python
		".pyi",         // Python
		".pyx",         // Cython
		".qml",         // QML
		".r",           // R
		".r2",          // Rebol
		".r2s",         // D
		".raku",        // Raku
		".rb",          // Ruby
		".rbs",         // Ruby
		".rbw",         // Ruby
		".re",          // Reason
		".rei",         // Reason
		".res",         // D
		".rexx",        // Rexx
		".rhtml",       // HTML
		".ring",        // Ring
		".rkt",         // Racket
		".rktd",        // Racket
		".rktl",        // Racket
		".rmd",         // R
		".robot",       // Robot Framework
		".rs",          // Rust
		".rsh",         // C
		".rss",         // RSS
		".rst",         // reStructuredText
		".rsvp",        // D
		".rt",          // RealTime
		".rtf",         // Rich Text Format
		".s",           // Assembly
		".sage",        // SageMath
		".sas",         // SAS
		".sass",        // Sass
		".scala",       // Scala
		".scm",         // Scheme
		".scss",        // Sass
		".sed",         // sed
		".self",        // Self
		".shader",      // Shader
		".sh",          // Shell script
		".shen",        // Shen
		".sig",         // Standard ML
		".sls",         // Scheme
		".sml",         // Standard ML
		".sol",         // Solidity
		".sqf",         // SQF
		".sql",         // SQL
		".ss",          // Scheme
		".st",          // Smalltalk
		".swift",       // Swift
		".t",           // Perl
		".tac",         // Python
		".tcc",         // C++
		".tcl",         // Tcl
		".tex",         // TeX
		".thy",         // Isabelle
		".toml",        // TOML
		".ts",          // TypeScript
		".tsx",         // TypeScript
		".tu",          // D
		".twig",        // Twig
		".uc",          // UnrealScript
		".ul",          // D
		".ur",          // Ur
		".urs",         // Ur
		".v",           // Verilog
		".vala",        // Vala
		".vapi",        // Vala
		".vb",          // Visual Basic
		".vba",         // Visual Basic for Applications
		".vbproj",      // Visual Basic
		".vbs",         // VBScript
		".vhd",         // VHDL
		".vhdl",        // VHDL
		".vim",         // Vim script
		".vue",         // Vue.js
		".w",           // W
		".w6",          // Perl 6
		".wat",         // WebAssembly
		".webidl",      // WebIDL
		".wisp",        // Wisp
		".wl",          // Wolfram Language
		".wsf",         // Windows Script File
		".wsgi",        // Python
		".wxs",         // XML
		".wxi",         // XML
		".wxl",         // XML
		".x",           // X
		".x10",         // X10
		".xht",         // XHTML
		".xhtml",       // XHTML
		".xi",          // X
		".xm",          // XML
		".xmi",         // XML
		".xpl",         // XProc
		".xq",          // XQuery
		".xql",         // XQuery
		".xqm",         // XQuery
		".xquery",      // XQuery
		".xqy",         // XQuery
		".xs",          // XS
		".xsl",         // XSLT
		".xslt",        // XSLT
		".xtend",       // Xtend
		".y",           // Yacc
		".yml",         // YAML
		".yaml",        // YAML
		".zeek",        // Zeek
		".zep",         // Zephir
		".zig",         // Zig
		".zsh",         // Zsh
	}
)

// ValidatePath validates the path
func ValidatePath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory checks if the path is a directory or not
func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

// IsSensitiveFile checks if the file or directory is sensitive or not
func IsSensitiveFile(path string) bool {
	name := ExtractName(path)

	// Check for sensitive files, if the starts with . then it is sensitive
	return name[0] == '.'
}

// IsUnwantedFilesAndFolders checks if the file or directory is unwanted or not
func IsUnwantedFilesAndFolders(path string) bool {
	// Check if the file or directory or not
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if the file or directory is unwanted
	if fileInfo.IsDir() && IsSensitiveFile(path) {
		return true
	}

	// Check if the file or directory is unwanted
	for _, unwantedFileAndFoler := range unwantedFilesAndFolders {
		if strings.Contains(path, unwantedFileAndFoler) {
			return true
		}
	}

	return false
}

// IsMarkdownFile checks if the file is a markdown file or not
func IsMarkdownFile(path string) bool {
	ext := filepath.Ext(path)

	return strings.Contains(ext, ".md")
}

// IsProgrammingFile checks if the file is a programming file or not
func IsProgrammingFile(path string) bool {
	ext := filepath.Ext(path)

	// Check if the file is a programming file, using binary search
	low := 0
	high := len(validProgrammingFileExtensions) - 1
	for low <= high {
		mid := low + (high-low)/2
		extension := validProgrammingFileExtensions[mid]

		if extension == ext {
			return true
		}

		if extension < ext {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return false
}

// IsSupportedFile checks if the file is valid or not
func IsSupportedFile(path string) bool {
	return !IsSensitiveFile(path) && !IsMarkdownFile(path) && IsProgrammingFile(path)
}
