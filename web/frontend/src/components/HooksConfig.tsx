import type { HookType, HooksConfig } from '../types';

interface HooksConfigProps {
  hooks: HooksConfig;
  onChange: (hooks: HooksConfig) => void;
}

const AVAILABLE_HOOKS: { id: HookType; name: string; description: string }[] = [
  { id: 'fmt', name: 'fmt', description: 'Format code with clang-format' },
  { id: 'lint', name: 'lint', description: 'Run clang-tidy static analysis' },
  { id: 'test', name: 'test', description: 'Run tests (blocking for pre-push)' },
  { id: 'flawfinder', name: 'flawfinder', description: 'Run Flawfinder security analysis' },
  { id: 'cppcheck', name: 'cppcheck', description: 'Run Cppcheck static analysis' },
  { id: 'check', name: 'check', description: 'Run code check' },
];

export function HooksConfig({ hooks, onChange }: HooksConfigProps) {
  const toggleHook = (hookType: HookType, hookStage: 'precommit' | 'prepush') => {
    const currentHooks = hooks[hookStage];
    const isSelected = currentHooks.includes(hookType);
    
    const newHooks = {
      ...hooks,
      [hookStage]: isSelected
        ? currentHooks.filter((h) => h !== hookType)
        : [...currentHooks, hookType],
    };
    
    onChange(newHooks);
  };

  const renderHookCheckbox = (hook: typeof AVAILABLE_HOOKS[0], stage: 'precommit' | 'prepush') => {
    const isSelected = hooks[stage].includes(hook.id);
    
    return (
      <label
        key={hook.id}
        className="flex items-center gap-1.5 cursor-pointer group p-1.5 rounded hover:bg-white/5 transition-colors"
      >
        <input
          type="checkbox"
          checked={isSelected}
          onChange={() => toggleHook(hook.id, stage)}
          className="checkbox-custom w-3.5 h-3.5"
        />
        <div className="flex-1 min-w-0">
          <span className={`text-xs font-mono transition-colors ${
            isSelected ? 'text-cyan-400' : 'text-gray-400 group-hover:text-gray-300'
          }`}>
            {hook.name}
          </span>
          <span className="text-[10px] text-gray-500 ml-1.5 block truncate" title={hook.description}>{hook.description}</span>
        </div>
      </label>
    );
  };

  return (
    <div className="card-glass rounded-2xl p-4 space-y-3">
      <h2 className="font-display font-semibold text-base text-white flex items-center gap-2">
        <svg className="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
        </svg>
        Git Hooks Configuration
      </h2>

      <div className="grid grid-cols-2 gap-3">
        {/* Pre-commit Hooks */}
        <div>
          <label className="block text-xs font-medium text-gray-400 mb-2">
            Pre-commit
            <span className="text-xs text-gray-500 ml-1">(before commit)</span>
          </label>
          <div className="bg-white/5 rounded-lg p-2 space-y-0.5 border border-white/10 max-h-[280px] overflow-y-auto">
            {AVAILABLE_HOOKS.map((hook) => renderHookCheckbox(hook, 'precommit'))}
          </div>
        </div>

        {/* Pre-push Hooks */}
        <div>
          <label className="block text-xs font-medium text-gray-400 mb-2">
            Pre-push
            <span className="text-xs text-gray-500 ml-1">(before push)</span>
          </label>
          <div className="bg-white/5 rounded-lg p-2 space-y-0.5 border border-white/10 max-h-[280px] overflow-y-auto">
            {AVAILABLE_HOOKS.map((hook) => renderHookCheckbox(hook, 'prepush'))}
          </div>
        </div>
      </div>
    </div>
  );
}

