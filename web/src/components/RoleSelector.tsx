import { Role } from '../types';
import { Send, Download } from 'lucide-react';

interface RoleSelectorProps {
  selectedRole: Role;
  onRoleChange: (role: Role) => void;
}

export default function RoleSelector({ selectedRole, onRoleChange }: RoleSelectorProps) {
  return (
    <div className="flex flex-col sm:flex-row gap-4 mb-8 w-full max-w-md">
      <button
        className={`flex-1 flex items-center justify-center gap-2 py-3 px-6 rounded-lg font-medium transition-all duration-200 ${
          selectedRole === 'receiver'
            ? 'bg-indigo-600 text-white shadow-lg scale-105'
            : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600'
        }`}
        onClick={() => onRoleChange('receiver')}
      >
        <Download size={20} />
        Receiver
      </button>
      
      <button
        className={`flex-1 flex items-center justify-center gap-2 py-3 px-6 rounded-lg font-medium transition-all duration-200 ${
          selectedRole === 'sender'
            ? 'bg-indigo-600 text-white shadow-lg scale-105'
            : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600'
        }`}
        onClick={() => onRoleChange('sender')}
      >
        <Send size={20} />
        Sender
      </button>
    </div>
  );
}