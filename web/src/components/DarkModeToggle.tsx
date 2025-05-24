import { Moon, Sun } from 'lucide-react';
import { useTheme } from '../context/ThemeContext';

export default function DarkModeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      onClick={toggleTheme}
      className="p-2 rounded-full transition-colors duration-200 ease-in-out
                 dark:bg-gray-800 dark:text-yellow-300 dark:hover:bg-gray-700
                 bg-gray-200 text-gray-700 hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500"
      aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {theme === 'dark' ? (
        <Sun size={20} className="transition-transform duration-300 ease-in-out hover:rotate-90" />
      ) : (
        <Moon size={20} className="transition-transform duration-300 ease-in-out hover:scale-110" />
      )}
    </button>
  );
}