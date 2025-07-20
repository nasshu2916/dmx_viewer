import { vi } from 'vitest'

// fetch APIのモック
global.fetch = vi.fn() as unknown as typeof fetch
