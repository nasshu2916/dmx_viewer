/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'dmx-dark-bg': '#1c1c1c',
        'dmx-medium-bg': '#282828',
        'dmx-light-bg': '#333333',
        'dmx-text-gray': '#9ca3af',
        'dmx-text-light': '#f0f0f0',
        'dmx-accent': '#007bff',
        'dmx-channel-active': '#28a745',
        'dmx-channel-inactive': '#444444',
        'dmx-border': '#444444',
      },
    },
  },
  plugins: [],
}
