import { memo } from 'react'

export type MobileTabKey = 'nodes' | 'viewer' | 'status'

type Props = {
  active: MobileTabKey
  onChange: (key: MobileTabKey) => void
}

const TABS: { key: MobileTabKey; label: string }[] = [
  { key: 'nodes', label: 'Nodes' },
  { key: 'viewer', label: 'Viewer' },
  { key: 'status', label: 'Status' },
]

function MobileTabsComponent({ active, onChange }: Props) {
  return (
    <nav className="md:hidden">
      <ul className="grid h-10 grid-cols-3 overflow-hidden rounded-md border border-dmx-border bg-dmx-medium-bg">
        {TABS.map(({ key, label }) => {
          const isActive = key === active
          return (
            <li className="contents" key={key}>
              <button
                aria-current={isActive ? 'page' : undefined}
                className={
                  'text-center text-sm font-medium transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-dmx-accent ' +
                  (isActive
                    ? 'bg-dmx-light-bg text-dmx-text-light'
                    : 'text-dmx-text-gray hover:text-dmx-text-light')
                }
                type="button"
                onClick={() => onChange(key)}
              >
                {label}
              </button>
            </li>
          )
        })}
      </ul>
    </nav>
  )
}

const MobileTabs = memo(MobileTabsComponent)

export default MobileTabs

