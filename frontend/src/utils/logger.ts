enum LogLevel {
  DEBUG,
  INFO,
  WARN,
  ERROR,
}

let currentLogLevel: LogLevel

const initializeLogger = () => {
  const logLevelStr = import.meta.env.VITE_LOG_LEVEL?.toUpperCase()
  switch (logLevelStr) {
    case 'DEBUG':
      currentLogLevel = LogLevel.DEBUG
      break
    case 'INFO':
      currentLogLevel = LogLevel.INFO
      break
    case 'WARN':
      currentLogLevel = LogLevel.WARN
      break
    case 'ERROR':
      currentLogLevel = LogLevel.ERROR
      break
    default:
      currentLogLevel = LogLevel.INFO // Default to INFO
  }
  console.log(`Logger initialized with level: ${LogLevel[currentLogLevel]}`)
}

// Initialize the logger when the module is loaded
initializeLogger()

const enableColors = import.meta.env.VITE_ENABLE_LOG_COLORS === 'true'

const COLORS = {
  DEBUG: '\x1b[36m', // Cyan
  INFO: '\x1b[32m', // Green
  WARN: '\x1b[33m', // Yellow
  ERROR: '\x1b[31m', // Red
  RESET: '\x1b[0m', // Reset color
}

const formatMessage = (level: keyof typeof COLORS, message: unknown[]) => {
  if (!enableColors) {
    return [`[${level}]`, ...message]
  }
  return [`${COLORS[level]}[${level}]${COLORS.RESET}`, ...message]
}

export const logger = {
  debug: (...args: unknown[]) => {
    if (currentLogLevel <= LogLevel.DEBUG) {
      console.debug(...formatMessage('DEBUG', args))
    }
  },
  info: (...args: unknown[]) => {
    if (currentLogLevel <= LogLevel.INFO) {
      console.info(...formatMessage('INFO', args))
    }
  },
  warn: (...args: unknown[]) => {
    if (currentLogLevel <= LogLevel.WARN) {
      console.warn(...formatMessage('WARN', args))
    }
  },
  error: (...args: unknown[]) => {
    if (currentLogLevel <= LogLevel.ERROR) {
      console.error(...formatMessage('ERROR', args))
    }
  },
}
