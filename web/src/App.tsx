import { useState } from 'react';
import { Role, FormData } from './types';
import { ThemeProvider, useTheme } from './context/ThemeContext';
import DarkModeToggle from './components/DarkModeToggle';
import RoleSelector from './components/RoleSelector';
import ReceiverForm from './components/ReceiverForm';
import SenderForm from './components/SenderForm';
import StatusDisplay from './components/StatusDisplay';
import { Share2 } from 'lucide-react';

function AppContent() {
  const [role, setRole] = useState<Role>('receiver');
  const [status, setStatus] = useState<'idle' | 'connecting' | 'connected' | 'error'>('idle');
  const [formData, setFormData] = useState<FormData | undefined>(undefined);
  const { theme } = useTheme();

  const handleSubmit = (data: FormData, success :boolean) => {
    setFormData(data);
    setStatus('connecting');
    
    setTimeout(() => {
      setStatus(success ? "connected" : "error");
    }, 2000);
  };

  return (
    <div className={`min-h-screen flex flex-col transition-colors duration-300 
                    ${theme === 'dark' ? 'bg-gray-900 text-white' : 'bg-gray-50 text-gray-900'}`}>
      <header className="py-4 px-6 flex justify-between items-center border-b dark:border-gray-800">
        <div className="flex items-center gap-2">
          <Share2 className="text-indigo-600 dark:text-indigo-400" size={28} />
          <h1 className="text-xl font-bold">Echo file transport</h1>
        </div>
        <DarkModeToggle />
      </header>

      <main className="flex-1 flex flex-col items-center justify-center px-4 py-10">
        <div className="w-full max-w-lg mx-auto flex flex-col items-center">
          <h2 className="text-2xl font-bold mb-6 text-center">Choose Your Role</h2>
          
          <RoleSelector selectedRole={role} onRoleChange={setRole} />
          
          {status === 'idle' && (
            <>
              <h3 className="text-lg font-medium mb-4 text-center">
                {role === 'receiver' 
                  ? 'Set up to receive a file' 
                  : 'Send a file to a receiver'}
              </h3>
              
              {role === 'receiver' 
                ? <ReceiverForm onSubmit={handleSubmit} /> 
                : <SenderForm onSubmit={handleSubmit} />
              }
            </>
          )}
          
          <StatusDisplay 
            role={role} 
            formData={formData} 
            status={status} 
            errorMessage={status === 'error' ? 'Failed to establish connection' : undefined}
          />
        </div>
      </main>

      <footer className="py-3 px-6 text-center text-sm text-gray-500 dark:text-gray-400 border-t dark:border-gray-800">
        <p>© 2025 Echo file transfer. All rights reserved.</p>
      </footer>
    </div>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  );
}

export default App;