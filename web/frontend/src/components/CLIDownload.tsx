import { useState, useEffect } from 'react';
import { fetchVersion, type VersionInfo } from '../api';

interface Architecture {
  id: string;
  name: string;
  filename: string;
}

const MACOS_ARCHS: Architecture[] = [
  { id: 'darwin-arm64', name: 'Apple Silicon', filename: 'cpx-darwin-arm64' },
  { id: 'darwin-amd64', name: 'Intel', filename: 'cpx-darwin-amd64' },
];

const LINUX_ARCHS: Architecture[] = [
  { id: 'linux-amd64', name: 'x86_64', filename: 'cpx-linux-amd64' },
  { id: 'linux-arm64', name: 'ARM64', filename: 'cpx-linux-arm64' },
];

export function CLIDownload() {
  const [versionInfo, setVersionInfo] = useState<VersionInfo>({
    version: '0.0.35',
    cli_version: '0.0.35',
    name: 'cpx',
    description: 'C++ Project Generator',
  });
  
  const [macArch, setMacArch] = useState(MACOS_ARCHS[0]);
  const [linuxArch, setLinuxArch] = useState(LINUX_ARCHS[0]);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    fetchVersion().then(setVersionInfo).catch(() => {});
  }, []);

  const INSTALL_SCRIPT = 'curl -f https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh';
  const BASE_URL = 'https://github.com/ozacod/cpx/releases/latest/download';

  const copyCommand = () => {
    navigator.clipboard.writeText(INSTALL_SCRIPT);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="max-w-4xl mx-auto px-4">
      {/* Header */}
      <div className="flex items-center gap-6 mb-6">
        <div className="w-20 h-20 rounded-2xl bg-gray-900/80 border border-white/10 flex items-center justify-center">
          <img src="/cpx.svg" alt="Cpx" className="w-14 h-14" />
        </div>
        <div>
          <div className="flex items-baseline gap-3">
            <h1 className="text-4xl font-light text-cyan-400 font-mono">{versionInfo.cli_version}</h1>
            <a 
              href="https://github.com/ozacod/cpx/releases" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-sm text-gray-400 hover:text-cyan-400"
            >
              Changelog →
            </a>
          </div>
          <p className="text-gray-500 text-sm">C++ Project Generator</p>
        </div>
      </div>

      {/* Quick Install */}
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <svg className="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span className="text-white text-sm font-medium">Quick Install (macOS & Linux)</span>
        </div>
        <div className="flex items-center gap-2 bg-gray-900/80 border border-white/10 rounded-lg px-3 py-2">
          <code className="flex-1 text-cyan-400 font-mono text-xs overflow-x-auto">
            {INSTALL_SCRIPT}
          </code>
          <button
            onClick={copyCommand}
            className="text-gray-400 hover:text-white p-1"
            title="Copy to clipboard"
          >
            {copied ? (
              <svg className="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            ) : (
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            )}
          </button>
        </div>
      </div>

      {/* Download Grid */}
      <div className="grid grid-cols-3 gap-4 mb-6">
        {/* macOS */}
        <div className="bg-gray-900/50 border border-white/10 rounded-xl p-4">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <svg className="w-5 h-5 text-gray-400" viewBox="0 0 24 24" fill="currentColor">
                <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
              </svg>
              <span className="text-white text-sm font-medium">macOS</span>
            </div>
            <select 
              value={macArch.id}
              onChange={(e) => setMacArch(MACOS_ARCHS.find(a => a.id === e.target.value) || MACOS_ARCHS[0])}
              className="bg-transparent text-gray-400 text-xs border-none outline-none cursor-pointer"
            >
              {MACOS_ARCHS.map(arch => (
                <option key={arch.id} value={arch.id} className="bg-gray-900">{arch.name}</option>
              ))}
            </select>
          </div>
          <a
            href={`${BASE_URL}/${macArch.filename}`}
            className="flex items-center justify-center gap-2 w-full py-2 bg-cyan-500 hover:bg-cyan-400 text-black text-sm font-semibold rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            Binary
          </a>
        </div>

        {/* Windows */}
        <div className="bg-gray-900/50 border border-white/10 rounded-xl p-4">
          <div className="flex items-center gap-2 mb-3">
            <svg className="w-5 h-5 text-gray-400" viewBox="0 0 24 24" fill="currentColor">
              <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
            </svg>
            <span className="text-white text-sm font-medium">Windows</span>
            <span className="text-gray-500 text-xs ml-auto">x64</span>
          </div>
          <a
            href={`${BASE_URL}/cpx-windows-amd64.exe`}
            className="flex items-center justify-center gap-2 w-full py-2 bg-cyan-500 hover:bg-cyan-400 text-black text-sm font-semibold rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            EXE
          </a>
        </div>

        {/* Linux */}
        <div className="bg-gray-900/50 border border-white/10 rounded-xl p-4">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <svg className="w-5 h-5 text-gray-400" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a.424.424 0 00-.11.135c-.26.268-.45.6-.663.839-.199.199-.485.267-.797.4-.313.136-.658.269-.864.68-.09.189-.136.394-.132.602 0 .199.027.4.055.536.058.399.116.728.04.97-.249.68-.28 1.145-.106 1.484.174.334.535.47.94.601.81.2 1.91.135 2.774.6.926.466 1.866.67 2.616.47.526-.116.97-.464 1.208-.946.587-.003 1.23-.269 2.26-.334.699-.058 1.574.267 2.577.2.025.134.063.198.114.333l.003.003c.391.778 1.113 1.132 1.884 1.071.771-.06 1.592-.536 2.257-1.306.631-.765 1.683-1.084 2.378-1.503.348-.199.629-.469.649-.853.023-.4-.2-.811-.714-1.376v-.097l-.003-.003c-.17-.2-.25-.535-.338-.926-.085-.401-.182-.786-.492-1.046h-.003c-.059-.054-.123-.067-.188-.135a.357.357 0 00-.19-.064c.431-1.278.264-2.55-.173-3.694-.533-1.41-1.465-2.638-2.175-3.483-.796-1.005-1.576-1.957-1.56-3.368.026-2.152.236-6.133-3.544-6.139z"/>
              </svg>
              <span className="text-white text-sm font-medium">Linux</span>
            </div>
            <select 
              value={linuxArch.id}
              onChange={(e) => setLinuxArch(LINUX_ARCHS.find(a => a.id === e.target.value) || LINUX_ARCHS[0])}
              className="bg-transparent text-gray-400 text-xs border-none outline-none cursor-pointer"
            >
              {LINUX_ARCHS.map(arch => (
                <option key={arch.id} value={arch.id} className="bg-gray-900">{arch.name}</option>
              ))}
            </select>
          </div>
          <a
            href={`${BASE_URL}/${linuxArch.filename}`}
            className="flex items-center justify-center gap-2 w-full py-2 bg-cyan-500 hover:bg-cyan-400 text-black text-sm font-semibold rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            Binary
          </a>
        </div>
      </div>


      {/* Footer */}
      <p className="text-xs text-gray-500 text-center">
        <a href="https://github.com/ozacod/cpx/blob/master/LICENSE" target="_blank" rel="noopener noreferrer" className="text-cyan-400 hover:text-cyan-300">MIT license</a>
        {' • '}
        <a href="https://github.com/ozacod/cpx" target="_blank" rel="noopener noreferrer" className="text-cyan-400 hover:text-cyan-300">GitHub</a>
      </p>
    </div>
  );
}
