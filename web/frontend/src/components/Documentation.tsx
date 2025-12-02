import { useState } from 'react';

type DocSection = 'getting-started' | 'commands' | 'configuration' | 'project-structure' | 'ci' | 'examples';
type DocSubSection = string;

interface DocItem {
  id: string;
  title: string;
  icon?: string;
  subsections?: { id: string; title: string }[];
}

export function Documentation() {
  const [activeSection, setActiveSection] = useState<DocSection>('getting-started');
  const [activeSubSection, setActiveSubSection] = useState<DocSubSection>('');
  const [expandedSections, setExpandedSections] = useState<Set<DocSection>>(new Set(['getting-started']));

  const docSections: Record<DocSection, DocItem> = {
    'getting-started': {
      id: 'getting-started',
      title: 'Getting Started',
      icon: '🚀',
      subsections: [
        { id: 'quick-start', title: 'Quick Start' },
        { id: 'installation', title: 'Installation' },
        { id: 'features', title: 'Features' },
      ],
    },
    'commands': {
      id: 'commands',
      title: 'Commands',
      icon: '⚡',
      subsections: [
        { id: 'project-management', title: 'Project Management' },
        { id: 'build-run', title: 'Build & Run' },
        { id: 'dependencies', title: 'Dependencies' },
        { id: 'code-quality', title: 'Code Quality' },
        { id: 'configuration', title: 'Configuration' },
        { id: 'other', title: 'Other' },
      ],
    },
    'configuration': {
      id: 'configuration',
      title: 'Configuration',
      icon: '⚙️',
      subsections: [
        { id: 'global-config', title: 'Global Configuration' },
        { id: 'project-config', title: 'Project Configuration' },
        { id: 'templates', title: 'Templates' },
        { id: 'hooks', title: 'Git Hooks' },
        { id: 'cmake-presets', title: 'CMake Presets' },
      ],
    },
    'project-structure': {
      id: 'project-structure',
      title: 'Project Structure',
      icon: '📁',
      subsections: [
        { id: 'generated-files', title: 'Generated Files' },
        { id: 'key-files', title: 'Key Files' },
      ],
    },
    'ci': {
      id: 'ci',
      title: 'CI/CD',
      icon: '🐳',
      subsections: [
        { id: 'cross-compilation', title: 'Cross-Compilation' },
        { id: 'ci-commands', title: 'CI Commands' },
        { id: 'setup', title: 'Setup' },
      ],
    },
    'examples': {
      id: 'examples',
      title: 'Examples',
      icon: '💡',
      subsections: [
        { id: 'creating-project', title: 'Creating a Project' },
        { id: 'adding-dependencies', title: 'Adding Dependencies' },
        { id: 'build-options', title: 'Build Options' },
        { id: 'testing', title: 'Testing' },
        { id: 'code-quality', title: 'Code Quality & Security' },
        { id: 'cross-compilation', title: 'Cross-Compilation' },
      ],
    },
  };

  const toggleSection = (section: DocSection) => {
    const newExpanded = new Set(expandedSections);
    if (newExpanded.has(section)) {
      newExpanded.delete(section);
    } else {
      newExpanded.add(section);
    }
    setExpandedSections(newExpanded);
  };

  const handleSectionClick = (section: DocSection) => {
    setActiveSection(section);
    if (!expandedSections.has(section)) {
      setExpandedSections(new Set([...expandedSections, section]));
    }
    setActiveSubSection(docSections[section].subsections?.[0]?.id || '');
  };

  const handleSubSectionClick = (section: DocSection, subSection: string) => {
    setActiveSection(section);
    setActiveSubSection(subSection);
    if (!expandedSections.has(section)) {
      setExpandedSections(new Set([...expandedSections, section]));
    }
  };

  return (
    <div className="animate-fade-in max-w-7xl mx-auto">
      <div className="flex gap-6">
        {/* Left Sidebar */}
        <aside className="w-64 flex-shrink-0">
          <nav className="sticky top-20 bg-black/20 backdrop-blur-sm border border-white/10 rounded-lg p-4">
            <div className="space-y-1">
              {Object.values(docSections).map((section) => {
                const isExpanded = expandedSections.has(section.id as DocSection);
                const hasSubsections = section.subsections && section.subsections.length > 0;
                
                return (
                  <div key={section.id} className="mb-1">
                    <button
                      onClick={() => hasSubsections ? toggleSection(section.id as DocSection) : handleSectionClick(section.id as DocSection)}
                      className={`w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm font-medium transition-all ${
                        activeSection === section.id
                          ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/30'
                          : 'text-gray-400 hover:text-white hover:bg-white/5'
                      }`}
                    >
                      <div className="flex items-center gap-2">
                        {section.icon && <span>{section.icon}</span>}
                        <span>{section.title}</span>
                      </div>
                      {hasSubsections && (
                        <svg
                          className={`w-4 h-4 transition-transform ${isExpanded ? 'rotate-90' : ''}`}
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                      )}
                    </button>
                    
                    {hasSubsections && isExpanded && (
                      <div className="ml-4 mt-1 space-y-1 border-l border-white/10 pl-3">
                        {section.subsections?.map((subsection) => (
                          <button
                            key={subsection.id}
                            onClick={() => handleSubSectionClick(section.id as DocSection, subsection.id)}
                            className={`w-full text-left px-3 py-1.5 rounded text-sm transition-all ${
                              activeSection === section.id && activeSubSection === subsection.id
                                ? 'text-cyan-400 bg-cyan-500/10'
                                : 'text-gray-500 hover:text-gray-300 hover:bg-white/5'
                            }`}
                          >
                            {subsection.title}
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </nav>
        </aside>

        {/* Main Content Area */}
        <main className="flex-1 min-w-0">
          {/* Breadcrumbs */}
          <div className="mb-6 flex items-center gap-2 text-sm text-gray-400">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            <span>/</span>
            <span className="px-2 py-0.5 bg-white/10 rounded text-white">
              {docSections[activeSection].title}
            </span>
            {activeSubSection && (
              <>
                <span>/</span>
                <span className="text-gray-300">{docSections[activeSection].subsections?.find(s => s.id === activeSubSection)?.title}</span>
              </>
            )}
          </div>

          {/* Content */}
          <div className="min-h-[600px]">
            {activeSection === 'getting-started' && <GettingStarted activeSubSection={activeSubSection} />}
            {activeSection === 'commands' && <Commands activeSubSection={activeSubSection} />}
            {activeSection === 'configuration' && <Configuration activeSubSection={activeSubSection} />}
            {activeSection === 'project-structure' && <ProjectStructure activeSubSection={activeSubSection} />}
            {activeSection === 'ci' && <CI activeSubSection={activeSubSection} />}
            {activeSection === 'examples' && <Examples activeSubSection={activeSubSection} />}
          </div>
        </main>
      </div>
    </div>
  );
}

function GettingStarted({ activeSubSection }: { activeSubSection: string }) {
  if (activeSubSection === 'installation') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Installation</h2>
          <div className="space-y-4">
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Quick Install</h3>
              <div className="bg-black/40 rounded-lg p-4 font-mono text-sm">
                <code className="text-cyan-400">curl -f https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh</code>
              </div>
              <p className="text-sm text-gray-400 mt-2">
                The installer will automatically detect your OS and architecture, download the latest cpx binary, 
                set up vcpkg, and configure cpx with the vcpkg root directory.
              </p>
            </div>

            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Manual Installation</h3>
              <div className="space-y-3">
                <div>
                  <p className="text-sm text-gray-400 mb-2">1. Download the binary for your platform:</p>
                  <div className="bg-black/40 rounded-lg p-3 font-mono text-xs text-cyan-400">
                    # Visit GitHub releases page<br />
                    https://github.com/ozacod/cpx/releases/latest
                  </div>
                </div>
                <div>
                  <p className="text-sm text-gray-400 mb-2">2. Make it executable and move to PATH:</p>
                  <div className="bg-black/40 rounded-lg p-3 font-mono text-xs text-cyan-400">
                    chmod +x cpx-linux-amd64<br />
                    sudo mv cpx-linux-amd64 /usr/local/bin/cpx
                  </div>
                </div>
                <div>
                  <p className="text-sm text-gray-400 mb-2">3. Configure vcpkg:</p>
                  <div className="bg-black/40 rounded-lg p-3 font-mono text-xs text-cyan-400">
                    cpx config set-vcpkg-root /path/to/vcpkg
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (activeSubSection === 'features') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Features</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">📦</span>
              <div>
                <h3 className="font-semibold text-white mb-1">vcpkg Integration</h3>
                <p className="text-sm text-gray-400">Direct integration with Microsoft vcpkg for dependency management</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">🔧</span>
              <div>
                <h3 className="font-semibold text-white mb-1">CMake Presets</h3>
                <p className="text-sm text-gray-400">Automatic CMakePresets.json generation for seamless IDE integration</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">✨</span>
              <div>
                <h3 className="font-semibold text-white mb-1">Code Quality Tools</h3>
                <p className="text-sm text-gray-400">Built-in clang-format, clang-tidy, Flawfinder, and Cppcheck integration</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">🔒</span>
              <div>
                <h3 className="font-semibold text-white mb-1">Security Analysis</h3>
                <p className="text-sm text-gray-400">Flawfinder and Cppcheck for security vulnerability detection</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">🧪</span>
              <div>
                <h3 className="font-semibold text-white mb-1">Sanitizers</h3>
                <p className="text-sm text-gray-400">AddressSanitizer, ThreadSanitizer, MemorySanitizer, and UBSan support</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">🧪</span>
              <div>
                <h3 className="font-semibold text-white mb-1">Testing Support</h3>
                <p className="text-sm text-gray-400">Automatic test framework setup (googletest, catch2, doctest)</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">🔄</span>
              <div>
                <h3 className="font-semibold text-white mb-1">CI/CD Integration</h3>
                <p className="text-sm text-gray-400">Generate GitHub Actions and GitLab CI workflows automatically</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">🐳</span>
              <div>
                <h3 className="font-semibold text-white mb-1">Cross-Compilation</h3>
                <p className="text-sm text-gray-400">Docker-based cross-compilation for multiple platforms</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">⚡</span>
              <div>
                <h3 className="font-semibold text-white mb-1">vcpkg Passthrough</h3>
                <p className="text-sm text-gray-400">All vcpkg commands work directly through cpx</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <span className="text-cyan-400 text-xl">📝</span>
              <div>
                <h3 className="font-semibold text-white mb-1">Configurable Git Hooks</h3>
                <p className="text-sm text-gray-400">Automatically install git hooks (pre-commit, pre-push) based on cpx.yaml configuration</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Default: Quick Start
  return (
    <div className="space-y-6">
      <div className="card-glass rounded-xl p-6">
        <h2 className="text-2xl font-bold text-white mb-4">Quick Start</h2>
        <div className="space-y-4">
          <div>
            <h3 className="text-lg font-semibold text-white mb-2">1. Install Cpx</h3>
            <div className="bg-black/40 rounded-lg p-4 font-mono text-sm">
              <code className="text-cyan-400">curl -f https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh</code>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-semibold text-white mb-2">2. Create a Project</h3>
            <div className="bg-black/40 rounded-lg p-4 font-mono text-sm space-y-2">
              <div>
                <p className="text-gray-500 text-xs mb-1"># Create executable project</p>
                <code className="text-cyan-400">cpx create my_app</code>
              </div>
              <div>
                <p className="text-gray-500 text-xs mb-1"># Create library project</p>
                <code className="text-cyan-400">cpx create my_lib --lib</code>
              </div>
              <div>
                <p className="text-gray-500 text-xs mb-1"># Create from template (default)</p>
                <code className="text-cyan-400">cpx create my_project --template default</code>
              </div>
              <div>
                <p className="text-gray-500 text-xs mb-1"># Create with Catch2 template</p>
                <code className="text-cyan-400">cpx create my_project --template catch</code>
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-semibold text-white mb-2">3. Build and Run</h3>
            <div className="bg-black/40 rounded-lg p-4 font-mono text-sm space-y-2">
              <div>
                <p className="text-gray-500 text-xs mb-1"># Navigate to project</p>
                <code className="text-cyan-400">cd my_app</code>
              </div>
              <div>
                <p className="text-gray-500 text-xs mb-1"># Build the project</p>
                <code className="text-cyan-400">cpx build</code>
              </div>
              <div>
                <p className="text-gray-500 text-xs mb-1"># Build and run</p>
                <code className="text-cyan-400">cpx run</code>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function Commands({ activeSubSection }: { activeSubSection: string }) {
  const allCommands = {
    'project-management': [
      { command: 'cpx create <name>', description: 'Create new project with default config', options: [
        { flag: '--template <name>', desc: 'Create project from template (default, catch, or path to .yaml file)' },
        { flag: '--lib', desc: 'Create library project' },
      ]},
    ],
    'build-run': [
      { command: 'cpx build', description: 'Compile the project (uses CMake presets if available)', options: [
        { flag: '--release', desc: 'Build in release mode' },
        { flag: '-O<level>', desc: 'Optimization level: 0, 1, 2, 3, s, fast' },
        { flag: '--clean', desc: 'Clean and rebuild' },
        { flag: '-j <n>', desc: 'Use n parallel jobs' },
      ]},
      { command: 'cpx run', description: 'Build and run executable', options: [
        { flag: '--release', desc: 'Run in release mode' },
      ]},
      { command: 'cpx test', description: 'Build and run tests', options: [
        { flag: '-v, --verbose', desc: 'Verbose test output' },
        { flag: '--filter <name>', desc: 'Filter tests by name' },
      ]},
      { command: 'cpx clean', description: 'Remove build artifacts', options: [
        { flag: '--all', desc: 'Also remove generated files' },
      ]},
    ],
    'dependencies': [
      { command: 'cpx add port <package>', description: 'Add dependency (calls: vcpkg add port <package>)' },
      { command: 'cpx remove <package>', description: 'Remove dependency (calls: vcpkg remove <package>)' },
      { command: 'cpx list', description: 'List installed packages (calls: vcpkg list)' },
      { command: 'cpx search <query>', description: 'Search packages (calls: vcpkg search)' },
    ],
    'code-quality': [
      { command: 'cpx fmt', description: 'Format code with clang-format', options: [
        { flag: '--check', desc: 'Check formatting without modifying' },
      ]},
      { command: 'cpx lint', description: 'Run clang-tidy static analysis', options: [
        { flag: '--fix', desc: 'Auto-fix lint issues' },
      ]},
        { flag: '--json', desc: 'Output results in JSON format' },
        { flag: '--quiet', desc: 'Only output findings' },
        { flag: '--strict', desc: 'Return non-zero exit code on findings' },
        { flag: '--no-git-ignore', desc: 'Scan untracked files (not just git-tracked files)' },
      ]},
      { command: 'cpx flawfinder', description: 'Run Flawfinder security analysis for C/C++', options: [
        { flag: '--minlevel <0-5>', desc: 'Minimum risk level to report (default: 1)' },
        { flag: '--html', desc: 'Output results in HTML format' },
        { flag: '--csv', desc: 'Output results in CSV format' },
        { flag: '--output <file>', desc: 'Output file path (required for HTML/CSV)' },
        { flag: '--dataflow', desc: 'Enable dataflow analysis' },
        { flag: '--quiet', desc: 'Quiet mode (minimal output)' },
        { flag: '--context <n>', desc: 'Number of lines of context to show (default: 2)' },
      ]},
      { command: 'cpx cppcheck', description: 'Run Cppcheck static analysis for C/C++', options: [
        { flag: '--enable <checks>', desc: 'Enable checks (all, style, performance, portability, etc.)' },
        { flag: '--xml', desc: 'Output results in XML format' },
        { flag: '--csv', desc: 'Output results in CSV format' },
        { flag: '--output <file>', desc: 'Output file path (for XML/CSV output)' },
        { flag: '--quiet', desc: 'Quiet mode (suppress progress messages)' },
        { flag: '--force', desc: 'Force checking of all configurations' },
        { flag: '--platform <name>', desc: 'Target platform (unix32, unix64, win32A, win64, etc.)' },
        { flag: '--std <standard>', desc: 'C/C++ standard (c++17, c++20, etc.)' },
      ]},
      { command: 'cpx check', description: 'Check code compiles with sanitizers', options: [
        { flag: '--asan', desc: 'Build with AddressSanitizer (detects memory errors)' },
        { flag: '--tsan', desc: 'Build with ThreadSanitizer (detects data races)' },
        { flag: '--msan', desc: 'Build with MemorySanitizer (detects uninitialized memory)' },
        { flag: '--ubsan', desc: 'Build with UndefinedBehaviorSanitizer (detects undefined behavior)' },
      ]},
    ],
    'configuration': [
      { command: 'cpx config set-vcpkg-root <path>', description: 'Set vcpkg installation directory' },
      { command: 'cpx config get-vcpkg-root', description: 'Get current vcpkg root' },
      { command: 'cpx hooks install', description: 'Install git hooks based on cpx.yaml configuration', options: [
        { flag: '', desc: 'Creates hooks for precommit/prepush if configured, .sample files otherwise' },
      ]},
    ],
    'other': [
      { command: 'cpx upgrade', description: 'Upgrade cpx to latest version' },
      { command: 'cpx version', description: 'Show version' },
      { command: 'cpx help', description: 'Show help' },
    ],
  };

  const section = activeSubSection || 'project-management';
  const commands = allCommands[section as keyof typeof allCommands] || [];

  return (
    <div className="space-y-6">
      <div className="card-glass rounded-xl p-6">
        <h2 className="text-2xl font-bold text-white mb-4">Commands</h2>
        <div className="space-y-3">
          {commands.map((cmd, i) => (
            <CommandItem key={i} {...cmd} />
          ))}
        </div>
        {section === 'dependencies' && (
          <div className="mt-4 p-3 bg-cyan-500/10 border border-cyan-500/20 rounded-lg">
            <p className="text-sm text-cyan-400">
              <strong>Note:</strong> All vcpkg commands pass through automatically. You can use{' '}
              <code className="bg-black/40 px-1 rounded">cpx install</code>,{' '}
              <code className="bg-black/40 px-1 rounded">cpx list</code>, etc.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}

function Configuration({ activeSubSection }: { activeSubSection: string }) {
  if (activeSubSection === 'global-config') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Global Configuration</h2>
          <p className="text-gray-400 mb-4">
            Cpx stores its configuration in:
          </p>
          <div className="bg-black/40 rounded-lg p-4 font-mono text-sm space-y-2">
            <div>
              <span className="text-gray-500"># Linux/macOS:</span>
              <code className="text-cyan-400 block mt-1">~/.config/cpx/config.yaml</code>
            </div>
            <div>
              <span className="text-gray-500"># Windows:</span>
              <code className="text-cyan-400 block mt-1">%APPDATA%/cpx/config.yaml</code>
            </div>
          </div>
          <div className="mt-4 bg-black/40 rounded-lg p-4">
            <p className="text-gray-500 text-xs mb-2">config.yaml:</p>
            <pre className="text-sm text-cyan-400">
{`vcpkg_root: "/path/to/vcpkg"`}
            </pre>
          </div>
        </div>
      </div>
    );
  }

  if (activeSubSection === 'project-config') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Project Configuration</h2>
          <p className="text-gray-400 mb-4">
            Dependencies are managed in <code className="text-cyan-400">vcpkg.json</code>, not cpx.yaml.
            The <code className="text-cyan-400">cpx.yaml</code> file is only used as a template for project creation.
          </p>
          
          <div className="bg-black/40 rounded-lg p-4 mb-4">
            <p className="text-gray-500 text-xs mb-2">vcpkg.json (auto-generated):</p>
            <pre className="text-sm text-cyan-400">
{`{
  "dependencies": [
    "spdlog",
    "fmt",
    "nlohmann-json"
  ]
}`}
            </pre>
          </div>

          <div className="bg-black/40 rounded-lg p-4">
            <p className="text-gray-500 text-xs mb-2">cpx.yaml (template only):</p>
            <pre className="text-sm text-cyan-400">
{`package:
  name: my_project
  version: "0.1.0"
  cpp_standard: 17
  project_type: exe

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest

hooks:
  precommit:
    - fmt
    - lint
  prepush:
    - test
`}
            </pre>
          </div>
        </div>
      </div>
    );
  }

  if (activeSubSection === 'templates') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Project Templates</h2>
          <p className="text-gray-400 mb-4">
            Cpx provides project templates that are automatically downloaded from the GitHub repository.
            Templates define the project structure, build configuration, testing framework, and git hooks.
          </p>

          <div className="space-y-6">
            <div>
              <h3 className="font-semibold text-white mb-3">Using Templates</h3>
              <div className="bg-black/40 rounded-lg p-4 mb-4">
                <p className="text-gray-500 text-xs mb-2">Create a project from a template:</p>
                <code className="text-cyan-400 block mb-2">cpx create my_project --template default</code>
                <code className="text-cyan-400 block">cpx create my_project --template catch</code>
              </div>
              <p className="text-sm text-gray-400 mb-4">
                If no template is specified, the <code className="text-cyan-400">default</code> template is automatically downloaded and used.
              </p>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-3">Available Templates</h3>
              
              <div className="space-y-4">
                <div className="bg-black/40 rounded-lg p-4">
                  <h4 className="font-semibold text-cyan-400 mb-2">default</h4>
                  <p className="text-sm text-gray-400 mb-3">
                    The default template uses Google Test framework and includes standard git hooks configuration.
                  </p>
                  <div className="bg-black/60 rounded p-3 mb-2">
                    <p className="text-gray-500 text-xs mb-1">Usage:</p>
                    <code className="text-cyan-400 text-sm">cpx create my_project --template default</code>
                  </div>
                  <div className="bg-black/60 rounded p-3">
                    <p className="text-gray-500 text-xs mb-1">Configuration:</p>
                    <pre className="text-xs text-cyan-400">
{`package:
  version: 0.1.0
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest

hooks:
  precommit:
    - fmt
    - lint
  prepush:
    - test
`}
                    </pre>
                  </div>
                </div>

                <div className="bg-black/40 rounded-lg p-4">
                  <h4 className="font-semibold text-cyan-400 mb-2">catch</h4>
                  <p className="text-sm text-gray-400 mb-3">
                    The catch template uses Catch2 test framework. Catch2 is automatically downloaded via FetchContent.
                  </p>
                  <div className="bg-black/60 rounded p-3 mb-2">
                    <p className="text-gray-500 text-xs mb-1">Usage:</p>
                    <code className="text-cyan-400 text-sm">cpx create my_project --template catch</code>
                  </div>
                  <div className="bg-black/60 rounded p-3">
                    <p className="text-gray-500 text-xs mb-1">Configuration:</p>
                    <pre className="text-xs text-cyan-400">
{`package:
  version: 0.1.0
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: catch2

hooks:
  precommit:
    - fmt
    - lint
  prepush:
    - test
`}
                    </pre>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-3">Template Features</h3>
              <ul className="text-sm text-gray-400 space-y-2 ml-4">
                <li>• <strong>Automatic Download:</strong> Templates are downloaded from GitHub when needed</li>
                <li>• <strong>No Local Storage:</strong> Templates are not stored locally, always fetched from the repository</li>
                <li>• <strong>Testing Framework:</strong> Choose between googletest (default) or catch2</li>
                <li>• <strong>Git Hooks:</strong> Templates can include pre-configured git hooks</li>
                <li>• <strong>Build Configuration:</strong> C++ standard, clang-format style, and library settings</li>
              </ul>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-3">Creating Custom Templates</h3>
              <p className="text-sm text-gray-400 mb-3">
                You can create your own templates by creating a YAML file with the same structure as the built-in templates.
                Templates should be placed in the <code className="text-cyan-400">templates/</code> directory of the cpx repository.
              </p>
              <div className="bg-black/40 rounded-lg p-4">
                <p className="text-gray-500 text-xs mb-2">Template structure:</p>
                <pre className="text-xs text-cyan-400">
{`package:
  version: 0.1.0
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest  # or catch2

hooks:
  precommit:
    - fmt
    - lint
  prepush:
    - test
`}
                </pre>
              </div>
              <p className="text-sm text-gray-400 mt-3">
                <strong>Note:</strong> The <code className="text-cyan-400">name</code> field is not included in templates as it is set from the project name during creation.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (activeSubSection === 'hooks') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Git Hooks Configuration</h2>
          <p className="text-gray-400 mb-4">
            Configure git hooks in <code className="text-cyan-400">cpx.yaml</code> to automatically run code quality checks.
            Hooks are automatically installed when creating a project from a template.
          </p>
          
          <div className="bg-black/40 rounded-lg p-4 mb-4">
            <p className="text-gray-500 text-xs mb-2">cpx.yaml hooks configuration:</p>
            <pre className="text-sm text-cyan-400">
{`hooks:
  precommit:
    - fmt      # Format code before commit
    - lint     # Run linter before commit
  prepush:
    - test     # Run tests before push
    - semgrep  # Run security checks before push`}
            </pre>
          </div>

          <div className="space-y-4">
            <div>
              <h3 className="font-semibold text-white mb-2">Supported Hook Checks</h3>
              <ul className="text-sm text-gray-400 space-y-1 ml-4">
                <li>• <code className="text-cyan-400">fmt</code> - Format code with clang-format</li>
                <li>• <code className="text-cyan-400">lint</code> - Run clang-tidy static analysis</li>
                <li>• <code className="text-cyan-400">test</code> - Run tests (blocking for pre-push)</li>
                <li>• <code className="text-cyan-400">flawfinder</code> - Run Flawfinder security analysis</li>
                <li>• <code className="text-cyan-400">cppcheck</code> - Run Cppcheck static analysis</li>
                <li>• <code className="text-cyan-400">check</code> - Run code check</li>
              </ul>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-2">Installation</h3>
              <p className="text-sm text-gray-400 mb-2">
                Hooks are automatically installed when creating a project from a template with hooks configured.
                You can also install them manually:
              </p>
              <div className="bg-black/40 rounded-lg p-3">
                <code className="text-cyan-400">cpx hooks install</code>
              </div>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-2">Behavior</h3>
              <ul className="text-sm text-gray-400 space-y-1 ml-4">
                <li>• <strong>Hooks configured in cpx.yaml</strong> → Creates actual hook files (e.g., <code className="text-cyan-400">pre-commit</code>)</li>
                <li>• <strong>Hooks NOT configured</strong> → Creates <code className="text-cyan-400">.sample</code> files (e.g., <code className="text-cyan-400">pre-commit.sample</code>)</li>
                <li>• <strong>No cpx.yaml</strong> → Uses defaults (fmt, lint for pre-commit; test for pre-push)</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Default: CMake Presets
  return (
    <div className="space-y-6">
      <div className="card-glass rounded-xl p-6">
        <h2 className="text-2xl font-bold text-white mb-4">CMake Presets</h2>
        <p className="text-gray-400 mb-4">
          Cpx generates <code className="text-cyan-400">CMakePresets.json</code> for seamless IDE integration:
        </p>
        <div className="space-y-4">
          <div>
            <h3 className="font-semibold text-white mb-2">CMakePresets.json</h3>
            <ul className="text-sm text-gray-400 space-y-1 ml-4">
              <li>• Uses environment variables (<code className="text-cyan-400">{'$env{VCPKG_ROOT}'}</code>)</li>
              <li>• Safe to commit to version control</li>
              <li>• <code className="text-cyan-400">VCPKG_ROOT</code> is automatically set by cpx build commands</li>
              <li>• Works seamlessly with IDEs like VS Code, CLion, and Qt Creator</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

function ProjectStructure({ activeSubSection }: { activeSubSection: string }) {
  if (activeSubSection === 'key-files') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Key Files</h2>
          <div className="space-y-4">
            <FileDescription
              name="vcpkg.json"
              description="vcpkg manifest file that lists all dependencies. Managed via 'vcpkg add port' commands."
            />
            <FileDescription
              name="CMakeLists.txt"
              description="Main CMake build file with vcpkg integration. Auto-generated by cpx create."
            />
            <FileDescription
              name="CMakePresets.json"
              description="CMake presets for IDE integration. Uses environment variables, safe to commit."
            />
            <FileDescription
              name="cpx.ci"
              description="Cross-compilation configuration for Docker-based builds. Created with empty targets."
            />
            <FileDescription
              name=".clang-format"
              description="Code formatting configuration. Auto-generated based on your style preference."
            />
          </div>
        </div>
      </div>
    );
  }

  // Default: Generated Files
  return (
    <div className="space-y-6">
      <div className="card-glass rounded-xl p-6">
        <h2 className="text-2xl font-bold text-white mb-4">Generated Project Structure</h2>
        <div className="bg-black/40 rounded-lg p-4 font-mono text-sm text-gray-300">
          <pre className="whitespace-pre-wrap">{`my_project/
├── CMakeLists.txt           # Main CMake file with vcpkg integration
├── CMakePresets.json        # CMake presets (safe to commit, uses env vars)
├── vcpkg.json               # vcpkg manifest (dependencies)
├── vcpkg-configuration.json # vcpkg configuration (auto-generated)
├── cpx.ci                 # Cross-compilation configuration
├── include/
│   └── my_project/
│       ├── my_project.hpp
│       └── version.hpp
├── src/
│   ├── main.cpp             # Main executable (if exe)
│   └── my_project.cpp
├── tests/
│   ├── CMakeLists.txt
│   └── test_main.cpp
├── .gitignore
├── .clang-format
└── README.md`}</pre>
        </div>
      </div>
    </div>
  );
}

function CI({ activeSubSection }: { activeSubSection: string }) {
  if (activeSubSection === 'ci-commands') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">CI Commands</h2>
          <div className="space-y-3">
            <CommandItem
              command="cpx ci"
              description="Build for all targets defined in cpx.ci"
              options={[
                { flag: '--target <name>', desc: 'Build only specific target' },
                { flag: '--rebuild', desc: 'Rebuild Docker images even if they exist' },
              ]}
            />
            <CommandItem
              command="cpx ci init --github-actions"
              description="Generate GitHub Actions workflow file (.github/workflows/ci.yml)"
            />
            <CommandItem
              command="cpx ci init --gitlab"
              description="Generate GitLab CI configuration file (.gitlab-ci.yml)"
            />
          </div>
        </div>
      </div>
    );
  }

  if (activeSubSection === 'setup') {
    return (
      <div className="space-y-6">
        <div className="card-glass rounded-xl p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Setup</h2>
          <div className="space-y-4">
            <div>
              <h3 className="font-semibold text-white mb-2">1. Download Dockerfiles</h3>
              <div className="bg-black/40 rounded-lg p-3 font-mono text-sm">
                <code className="text-cyan-400">cpx upgrade</code>
              </div>
              <p className="text-sm text-gray-400 mt-2">
                This downloads Dockerfiles to <code className="text-cyan-400">~/.config/cpx/dockerfiles/</code>
              </p>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-2">2. Configure cpx.ci</h3>
              <p className="text-sm text-gray-400 mb-2">
                Edit <code className="text-cyan-400">cpx.ci</code> in your project root and add targets:
              </p>
              <div className="bg-black/40 rounded-lg p-3 font-mono text-sm">
                <code className="text-cyan-400"># cpx.ci is created automatically with empty targets</code>
              </div>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-2">3. Generate CI Workflows (Optional)</h3>
              <div className="bg-black/40 rounded-lg p-3 font-mono text-sm space-y-2">
                <div>
                  <p className="text-gray-500 text-xs mb-1"># Generate GitHub Actions workflow</p>
                  <code className="text-cyan-400">cpx ci init --github-actions</code>
                </div>
                <div>
                  <p className="text-gray-500 text-xs mb-1"># Generate GitLab CI configuration</p>
                  <code className="text-cyan-400">cpx ci init --gitlab</code>
                </div>
              </div>
              <p className="text-sm text-gray-400 mt-2">
                These commands create CI workflow files that call <code className="text-cyan-400">cpx ci</code> automatically.
              </p>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-2">4. Build for Multiple Platforms</h3>
              <div className="bg-black/40 rounded-lg p-3 font-mono text-sm">
                <code className="text-cyan-400">cpx ci</code>
              </div>
              <p className="text-sm text-gray-400 mt-2">
                Artifacts will be in the <code className="text-cyan-400">out/</code> directory
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Default: Cross-Compilation
  return (
    <div className="space-y-6">
      <div className="card-glass rounded-xl p-6">
        <h2 className="text-2xl font-bold text-white mb-4">Cross-Compilation with Docker</h2>
        <p className="text-gray-400 mb-4">
          Cpx supports cross-compilation using Docker containers. The <code className="text-cyan-400">cpx.ci</code> file
          is automatically created when you run <code className="text-cyan-400">cpx create</code>.
        </p>

        <div className="bg-black/40 rounded-lg p-4 mb-4">
          <p className="text-gray-500 text-xs mb-2">cpx.ci example:</p>
          <pre className="text-sm text-cyan-400">
{`targets:
  - name: linux-amd64
    dockerfile: Dockerfile.linux-amd64
    image: cpx-linux-amd64
    triplet: x64-linux
    platform: linux/amd64

build:
  type: Release
  optimization: 2
  jobs: 0

output: out`}
          </pre>
        </div>
      </div>
    </div>
  );
}

function Examples({ activeSubSection }: { activeSubSection: string }) {
  const examples = {
    'creating-project': [
      { title: 'Basic Executable', code: `# Create a simple executable project\ncpx create my_app\ncd my_app\ncpx build\ncpx run` },
      { title: 'Library Project', code: `# Create a library project\ncpx create my_lib --lib\ncd my_lib\ncpx build` },
      { title: 'From Template', code: `# Create from default template (googletest)\ncpx create my_project --template default\n\n# Create with Catch2 template\ncpx create my_project --template catch\n\ncd my_project\ncpx build` },
    ],
    'adding-dependencies': [
      { title: 'Using vcpkg Commands', code: `# Add dependencies directly\ncpx add port spdlog\ncpx add port fmt\ncpx add port nlohmann-json\n\n# Or use vcpkg commands directly\ncpx install spdlog\ncpx list` },
      { title: 'Manual vcpkg.json Edit', code: `# Edit vcpkg.json directly\n{\n  "dependencies": [\n    "spdlog",\n    "fmt",\n    "nlohmann-json"\n  ]\n}` },
    ],
    'build-options': [
      { title: 'Release Build', code: `cpx build --release` },
      { title: 'Optimization Levels', code: `cpx build -O3        # Maximum optimization\ncpx build -O2        # Standard release (default)\ncpx build -O1        # Light optimization\ncpx build -O0        # No optimization (debug)` },
      { title: 'Parallel Build', code: `cpx build -j 8       # Use 8 parallel jobs\ncpx build -j 4       # Use 4 parallel jobs` },
    ],
    'testing': [
      { title: 'Run All Tests', code: `cpx test` },
      { title: 'Verbose Test Output', code: `cpx test -v` },
      { title: 'Filter Tests', code: `cpx test --filter MyTest` },
    ],
    'code-quality': [
      { title: 'Configure Git Hooks', code: `# Add to cpx.yaml\nhooks:\n  precommit:\n    - fmt\n    - lint\n  prepush:\n    - test\n\n# Install hooks (auto-installed on project creation)\ncpx hooks install` },
      { title: 'Flawfinder Analysis', code: `# Basic scan\ncpx flawfinder\n\n# HTML report\ncpx flawfinder --html --output report.html\n\n# CSV output with dataflow analysis\ncpx flawfinder --csv --output report.csv --dataflow` },
      { title: 'Cppcheck Static Analysis', code: `# Full analysis\ncpx cppcheck\n\n# XML report\ncpx cppcheck --xml --output report.xml\n\n# Specific checks only\ncpx cppcheck --enable style,performance` },
      { title: 'Sanitizer Checks', code: `# AddressSanitizer (memory errors)\ncpx check --asan\n\n# ThreadSanitizer (data races)\ncpx check --tsan\n\n# MemorySanitizer (uninitialized memory)\ncpx check --msan\n\n# UndefinedBehaviorSanitizer\ncpx check --ubsan` },
    ],
    'cross-compilation': [
      { title: 'Build for All Targets', code: `# Build for all targets in cpx.ci\ncpx ci` },
      { title: 'Build Specific Target', code: `# Build only for linux-amd64\ncpx ci --target linux-amd64` },
      { title: 'Rebuild Docker Images', code: `# Force rebuild of Docker images\ncpx ci --rebuild` },
      { title: 'Generate GitHub Actions Workflow', code: `# Create .github/workflows/ci.yml\ncpx ci init --github-actions` },
      { title: 'Generate GitLab CI Configuration', code: `# Create .gitlab-ci.yml\ncpx ci init --gitlab` },
    ],
  };

  const section = activeSubSection || 'creating-project';
  const sectionExamples = examples[section as keyof typeof examples] || [];

  return (
    <div className="space-y-6">
      <div className="card-glass rounded-xl p-6">
        <h2 className="text-2xl font-bold text-white mb-4">Examples</h2>
        <div className="space-y-4">
          {sectionExamples.map((example, i) => (
            <ExampleSection key={i} title={example.title} code={example.code} />
          ))}
        </div>
      </div>
    </div>
  );
}

function CommandItem({ 
  command, 
  description, 
  options 
}: { 
  command: string; 
  description: string; 
  options?: { flag: string; desc: string }[] 
}) {
  return (
    <div className="border-b border-white/5 pb-3 last:border-0 last:pb-0">
      <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-2">
        <div className="flex-1">
          <code className="text-cyan-400 font-mono text-sm">{command}</code>
          <p className="text-gray-400 text-sm mt-1">{description}</p>
        </div>
      </div>
      {options && options.length > 0 && (
        <div className="mt-2 ml-4 space-y-1">
          {options.map((opt, i) => (
            <div key={i} className="text-sm">
              <code className="text-cyan-400/80">{opt.flag}</code>
              <span className="text-gray-500 ml-2">{opt.desc}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function FileDescription({ name, description }: { name: string; description: string }) {
  return (
    <div className="border-b border-white/5 pb-3 last:border-0 last:pb-0">
      <code className="text-cyan-400 font-mono text-sm">{name}</code>
      <p className="text-gray-400 text-sm mt-1">{description}</p>
    </div>
  );
}

function ExampleSection({ title, code }: { title: string; code: string }) {
  return (
    <div>
      <h3 className="font-semibold text-white mb-2">{title}</h3>
      <div className="bg-black/40 rounded-lg p-4 font-mono text-sm">
        <pre className="text-cyan-400 whitespace-pre-wrap">{code}</pre>
      </div>
    </div>
  );
}
