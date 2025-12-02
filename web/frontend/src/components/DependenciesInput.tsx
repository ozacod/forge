import { useState } from 'react';

interface DependenciesInputProps {
  dependencies: string[];
  onChange: (dependencies: string[]) => void;
}

export function DependenciesInput({ dependencies, onChange }: DependenciesInputProps) {
  const [inputValue, setInputValue] = useState('');

  const handleAdd = () => {
    const trimmed = inputValue.trim();
    if (trimmed && !dependencies.includes(trimmed)) {
      onChange([...dependencies, trimmed]);
      setInputValue('');
    }
  };

  const handleRemove = (dep: string) => {
    onChange(dependencies.filter(d => d !== dep));
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAdd();
    }
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const pasted = e.clipboardData.getData('text');
    const lines = pasted.split('\n').map(l => l.trim()).filter(l => l);
    const newDeps = [...dependencies];
    lines.forEach(line => {
      if (line && !newDeps.includes(line)) {
        newDeps.push(line);
      }
    });
    onChange(newDeps);
  };

  return (
    <div className="card-glass rounded-2xl p-4 space-y-3">
      <h2 className="font-display font-semibold text-base text-white flex items-center gap-2">
        <svg className="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        Dependencies (vcpkg packages)
        <span className="ml-auto text-xs font-mono text-cyan-400">
          {dependencies.length}
        </span>
      </h2>

      <div className="flex gap-2">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
          onPaste={handlePaste}
          placeholder="e.g., spdlog, fmt, nlohmann-json"
          className="input-field flex-1 px-3 py-2 rounded-lg text-white font-mono text-sm"
        />
        <button
          onClick={handleAdd}
          disabled={!inputValue.trim() || dependencies.includes(inputValue.trim())}
          className="bg-cyan-500 hover:bg-cyan-400 disabled:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed text-black px-3 py-2 rounded-lg font-semibold text-xs transition-colors"
        >
          Add
        </button>
      </div>

      {dependencies.length > 0 ? (
        <div className="space-y-1.5 max-h-[200px] overflow-y-auto pr-2">
          {dependencies.map((dep) => (
            <div
              key={dep}
              className="flex items-center justify-between bg-white/5 rounded-lg px-2.5 py-1.5 group"
            >
              <div className="flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-cyan-400" />
                <span className="text-xs text-white font-mono">{dep}</span>
              </div>
              <button
                onClick={() => handleRemove(dep)}
                className="opacity-0 group-hover:opacity-100 p-1 text-gray-500 hover:text-red-400 transition-all"
                title="Remove"
              >
                <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-4 text-gray-500">
          <p className="text-xs">No dependencies</p>
          <p className="text-xs mt-0.5 text-gray-600">Add vcpkg package names above</p>
        </div>
      )}

      <div className="text-xs text-gray-500 pt-1.5 border-t border-white/10">
        <p>💡 Enter vcpkg package names or paste multiple (one per line)</p>
      </div>
    </div>
  );
}

