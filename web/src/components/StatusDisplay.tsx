import React from 'react';
import { FormData, Role } from '../types';
import { CheckCircle2, Clock, AlertCircle } from 'lucide-react';

type Status = 'idle' | 'connecting' | 'connected' | 'error';

interface StatusDisplayProps {
  role: Role;
  formData?: FormData;
  status: Status;
  errorMessage?: string;
}

export default function StatusDisplay({ role, formData, status, errorMessage }: StatusDisplayProps) {
  if (status === 'idle') return null;
  
  const getStatusIcon = () => {
    switch (status) {
      case 'connecting':
        return <Clock className="text-yellow-500 animate-pulse" size={24} />;
      case 'connected':
        return <CheckCircle2 className="text-green-500" size={24} />;
      case 'error':
        return <AlertCircle className="text-red-500" size={24} />;
      default:
        return null;
    }
  };

  const getStatusText = () => {
    switch (status) {
      case 'connecting':
        return role === 'receiver' ? 'Waiting for connection...' : 'Connecting to receiver...';
      case 'connected':
        return role === 'receiver' ? 'Connected and ready to receive' : 'Connected and ready to send';
      case 'error':
        return errorMessage || 'An error occurred';
      default:
        return '';
    }
  };

  return (
    <div className="mt-8 w-full max-w-md">
      <div className={`p-4 rounded-lg flex items-center gap-3 
                      ${status === 'connecting' ? 'bg-yellow-100 dark:bg-yellow-900/30' : ''}
                      ${status === 'connected' ? 'bg-green-100 dark:bg-green-900/30' : ''}
                      ${status === 'error' ? 'bg-red-100 dark:bg-red-900/30' : ''}
                      transition-colors duration-300`}>
        {getStatusIcon()}
        <div className="flex-1">
          <p className="font-medium dark:text-white">{getStatusText()}</p>
          {formData && status !== 'error' && (
            <div className="mt-2 text-sm text-gray-600 dark:text-gray-300">
              <p>Local Port: {formData.localPort}</p>
              <p>Remote Address: {formData.remoteAddr}</p>
              {role === 'sender' && formData.filePath && (
                <p>File: {formData.filePath}</p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}