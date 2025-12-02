import type { ClangFormatStyle, TestingFramework } from '../types';

interface ProjectConfigProps {
  cppStandard: number;
  onCppStandardChange: (standard: number) => void;
  includeTests: boolean;
  onIncludeTestsChange: (include: boolean) => void;
  testingFramework: TestingFramework;
  onTestingFrameworkChange: (framework: TestingFramework) => void;
  buildShared: boolean;
  onBuildSharedChange: (shared: boolean) => void;
  clangFormatStyle: ClangFormatStyle;
  onClangFormatStyleChange: (style: ClangFormatStyle) => void;
}

const CPP_STANDARDS = [11, 14, 17, 20, 23];

const CLANG_FORMAT_STYLES: { id: ClangFormatStyle; name: string; description: string }[] = [
  { id: 'Google', name: 'Google', description: 'Google C++ style guide' },
  { id: 'LLVM', name: 'LLVM', description: 'LLVM coding standards' },
  { id: 'Chromium', name: 'Chromium', description: 'Chromium project style' },
  { id: 'Mozilla', name: 'Mozilla', description: 'Mozilla coding style' },
  { id: 'WebKit', name: 'WebKit', description: 'WebKit coding style' },
  { id: 'Microsoft', name: 'Microsoft', description: 'Microsoft C++ style' },
  { id: 'GNU', name: 'GNU', description: 'GNU coding standards' },
];

const TESTING_FRAMEWORKS: { id: TestingFramework; name: string; description: string }[] = [
  { id: 'googletest', name: 'GoogleTest', description: 'Google\'s C++ testing framework with mocking' },
  { id: 'catch2', name: 'Catch2', description: 'Modern C++ test framework with BDD support' },
  { id: 'doctest', name: 'doctest', description: 'Fast single-header testing framework' },
  { id: 'none', name: 'None', description: 'No testing framework' },
];

export function ProjectConfig({
  cppStandard,
  onCppStandardChange,
  includeTests,
  onIncludeTestsChange,
  testingFramework,
  onTestingFrameworkChange,
  buildShared,
  onBuildSharedChange,
  clangFormatStyle,
  onClangFormatStyleChange,
}: ProjectConfigProps) {

  return (
    <div className="card-glass rounded-2xl p-4 space-y-3">
      <h2 className="font-display font-semibold text-base text-white flex items-center gap-2">
        <svg className="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        Project Configuration
      </h2>

      <div className="space-y-3">
        <div>
          <label className="block text-xs font-medium text-gray-400 mb-1.5">
            C++ Standard
          </label>
          <div className="grid grid-cols-5 gap-1.5">
            {CPP_STANDARDS.map((std) => (
              <button
                key={std}
                onClick={() => onCppStandardChange(std)}
                className={`py-1.5 rounded-lg font-mono text-xs transition-all ${
                  cppStandard === std
                    ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40'
                    : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                }`}
              >
                C++{std}
              </button>
            ))}
          </div>
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div>
            <label className="block text-xs font-medium text-gray-400 mb-1.5">
              Testing Framework
            </label>
            <div className="grid grid-cols-2 gap-1.5">
              {TESTING_FRAMEWORKS.map((fw) => (
                <button
                  key={fw.id}
                  onClick={() => {
                    onTestingFrameworkChange(fw.id);
                    if (fw.id !== 'none' && !includeTests) {
                      onIncludeTestsChange(true);
                    }
                    if (fw.id === 'none' && includeTests) {
                      onIncludeTestsChange(false);
                    }
                  }}
                  title={fw.description}
                  className={`py-1.5 px-2 rounded-lg font-mono text-xs transition-all ${
                    testingFramework === fw.id
                      ? 'bg-green-500/20 text-green-400 border border-green-500/40'
                      : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                  }`}
                >
                  {fw.name}
                </button>
              ))}
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-gray-400 mb-1.5">
              Clang-Format Style
            </label>
            <div className="grid grid-cols-2 gap-1.5">
              {CLANG_FORMAT_STYLES.slice(0, 4).map((style) => (
                <button
                  key={style.id}
                  onClick={() => onClangFormatStyleChange(style.id)}
                  title={style.description}
                  className={`py-1.5 px-2 rounded-lg font-mono text-xs transition-all ${
                    clangFormatStyle === style.id
                      ? 'bg-purple-500/20 text-purple-400 border border-purple-500/40'
                      : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                  }`}
                >
                  {style.name}
                </button>
              ))}
            </div>
            <div className="grid grid-cols-3 gap-1.5 mt-1.5">
              {CLANG_FORMAT_STYLES.slice(4).map((style) => (
                <button
                  key={style.id}
                  onClick={() => onClangFormatStyleChange(style.id)}
                  title={style.description}
                  className={`py-1.5 px-2 rounded-lg font-mono text-xs transition-all ${
                    clangFormatStyle === style.id
                      ? 'bg-purple-500/20 text-purple-400 border border-purple-500/40'
                      : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                  }`}
                >
                  {style.name}
                </button>
              ))}
            </div>
          </div>
        </div>

        <div>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={buildShared}
              onChange={(e) => onBuildSharedChange(e.target.checked)}
              className="checkbox-custom w-4 h-4"
            />
            <div>
              <span className="text-xs font-medium text-gray-400">Build Shared Libraries</span>
              <span className="text-xs text-gray-500 block">Shared (.so/.dylib/.dll) instead of static</span>
            </div>
          </label>
        </div>
      </div>
    </div>
  );
}
