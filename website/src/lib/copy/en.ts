export const en = {
  meta: {
    title: "NeoHM — Home Assistant in your menu bar",
    description:
      "A menu bar app for Home Assistant. Switch devices by room and control air conditioners from a panel that opens in the menu bar. Native builds for macOS, with Windows and Linux on the way.",
  },
  nav: {
    features: "Features",
    screens: "A closer look",
    setup: "Install",
    faq: "FAQ",
    github: "GitHub",
    download: "Download",
    language: "Language",
  },
  hero: {
    badge: "A menu bar app for Home Assistant",
    titleLead: "Your home,",
    titleAccent: "in the menu bar.",
    body: "NeoHM puts every room of your Home Assistant instance one click away. Switch lights, open curtains, tune the AC — without breaking your flow.",
    primary: "Download NeoHM",
    secondary: "View on GitHub",
    meta: ["macOS · Windows · Linux", "All architectures", "5 languages"],
    note: "Free and open source. macOS today; Windows and Linux on the way.",
  },
  features: {
    eyebrow: "Why NeoHM",
    title: "Everything at a glance. Nothing in the way.",
    body: "A panel that shows the state of your home, with controls that go as deep as each device allows.",
    items: [
      {
        title: "Rooms first",
        body: "Devices appear as tiles under their area, with their state and a warm glow while they're on. Grouping comes from your Home Assistant registries.",
      },
      {
        title: "Click to switch",
        body: "Click a tile to turn its device on or off; a device without a switch opens its page instead. The arrow always goes deeper, to dimming, colors, moving, playing or locking.",
      },
      {
        title: "Climate in full",
        body: "Power, temperature, mode, fan speed, airflow, Eco Mode and Sleep, shown as far as the unit supports each one.",
      },
      {
        title: "A quiet menu bar",
        body: "The icon stays the same whether devices are on or off, with nothing beside it. Click it when you need the panel, and click again to close it.",
      },
      {
        title: "Settings, inside the panel",
        body: "Language, address and token stay a page of the panel, never another window. Choose English, 简体中文, 繁體中文, 日本語 or 한국어.",
      },
      {
        title: "Every device type",
        body: "Lights, plugs, fans, covers, locks, media players and sensors each get the level, colors, actions and readings they actually support.",
      },
    ],
  },
  showcase: {
    eyebrow: "A closer look",
    title: "The whole home, one click deep.",
    body: "The panel follows the shape of your home. These are real screenshots of the app, not mockups.",
    tabs: {
      home: {
        label: "Home",
        title: "The house, by room",
        body: "Every device as a tile under its area. Lit tiles are on, and the line beneath each name carries the state that matters.",
        points: [
          "Grouped by area from Home Assistant",
          "State, brightness, position or temperature",
          "The arrow opens the device page",
        ],
      },
      climate: {
        label: "Climate",
        title: "A real climate remote",
        body: "The room temperature sits above the set point. Below it, a slider and the controls your unit offers — nothing it doesn't.",
        points: [
          "Target temperature in 1° steps",
          "Mode, fan speed, airflow, Eco and Sleep",
          "Room temperature under the name",
        ],
      },
      devices: {
        label: "Devices",
        title: "A page for every device",
        body: "A light offers brightness and color, a cover offers position, a TV offers playback, a lock offers its latch. Each page is shaped by the device.",
        points: [
          "Brightness, color, position, speed or volume",
          "Open, close, lock, unlock, play or pause",
          "Color swatches, or any color from a picker",
          "Quarter stops for a cover or valve",
        ],
        shots: ["Light", "Curtain", "TV"],
      },
      settings: {
        label: "Settings",
        title: "Connect once, then forget it",
        body: "One address, one token, one language. Settings is a page of the panel, opened from the gear, so setup never becomes another window.",
        points: [
          "A long-lived token and your address",
          "Language follows your system by default",
          "Everything stays on your computer",
        ],
      },
    },
  },
  languages: {
    eyebrow: "Languages",
    title: "Five languages, one panel.",
    body: "NeoHM follows the language of your computer, or you can choose one in Settings.",
  },
  setup: {
    eyebrow: "Install",
    title: "From download to your menu bar.",
    body: "NeoHM is available for macOS 13 or later, with Windows and Linux on the way. Install it, create an access token, and your home is one click away.",
    steps: [
      {
        title: "Download NeoHM",
        body: "Get the latest macOS release from GitHub and drag NeoHM into your Applications folder.",
      },
      {
        title: "Create an access token",
        body: "In your Home Assistant profile, create a long-lived access token.",
      },
      {
        title: "Connect, then close Settings",
        body: "Open NeoHM, enter the address and token, and connect. Your rooms are waiting.",
      },
    ],
    downloadTitle: "NeoHM for macOS",
    downloadBody: "Latest release · macOS 13 or later",
    download: "Download .dmg",
    releases: "All releases",
    free: "Free and open source",
    platforms: [
      { name: "macOS", status: "Available now" },
      { name: "Windows", status: "Coming soon" },
      { name: "Linux", status: "Coming soon" },
    ],
    architectures: "Apple silicon · Intel · x64 · ARM64",
    privacyTitle: "Local by design",
    privacyBody:
      "The address, token and language are saved in ~/.neohm/config.json, readable only by your user. NeoHM talks to your Home Assistant instance and nowhere else.",
    requirements: [
      "macOS 13+",
      "Windows & Linux next",
      "All architectures",
      "Home Assistant",
    ],
  },
  faq: {
    eyebrow: "FAQ",
    title: "Questions, answered.",
    body: "The short version: one Home Assistant, one computer, nothing in between.",
    more: "Still have a question?",
    link: "Open an issue on GitHub",
    items: [
      {
        question: "Is NeoHM free?",
        answer:
          "Yes. NeoHM is free and open source. Download the latest release from GitHub; the source is there too.",
      },
      {
        question: "Which platforms will it support?",
        answer:
          "NeoHM is available for macOS 13 or later today. Windows and Linux versions are on the way, with builds for every architecture: Apple silicon and Intel on macOS, x64 and ARM64 on Windows and Linux.",
      },
      {
        question: "Which Home Assistant setup does it need?",
        answer:
          "Any instance your computer can reach, on your local network or through a remote address. NeoHM reads device states over REST and loads the area and device registries over Home Assistant's WebSocket API.",
      },
      {
        question: "Where does my access token go?",
        answer:
          "Into ~/.neohm/config.json with the rest of your settings, readable only by your user. NeoHM sends it to the address you enter and nowhere else.",
      },
      {
        question: "Can I use more than one home?",
        answer:
          "NeoHM controls one Home Assistant instance at a time. The address and token in Settings decide which one.",
      },
      {
        question: "How quickly does it update?",
        answer:
          "Every 3 seconds while the panel is open. NeoHM reads nothing while the panel is closed, and a change you make shows at once.",
      },
      {
        question: "What if the menu bar icon can't be created?",
        answer:
          "NeoHM opens the same pages in a normal window instead, so you can still use it.",
      },
      {
        question: "Which languages does it speak?",
        answer:
          "English, Simplified Chinese, Traditional Chinese, Japanese and Korean. It follows your system by default.",
      },
    ],
  },
  cta: {
    title: ["Open the panel.", "Skip the app switching."],
    body: "NeoHM sends every change straight to Home Assistant. There is nothing else to keep in sync.",
    primary: "Download NeoHM",
    secondary: "Read the source",
  },
  footer: {
    tagline: "A menu bar for Home Assistant.",
    note: "NeoHM is an independent project and is not affiliated with Home Assistant or Nabu Casa.",
    links: {
      github: "GitHub",
      readme: "README",
      architecture: "Architecture",
      changelog: "Changelog",
    },
    project: "Project",
    builtWith: "Built with Go, MyGo and Onevra.",
    madeBy: "For macOS, Windows and Linux.",
  },
  alts: {
    home: "The NeoHM panel with Home Assistant devices grouped by room; the ceiling light and curtain are on.",
    climate:
      "An air conditioner page in NeoHM with a temperature slider, Eco Mode and mode controls.",
    settings:
      "NeoHM settings with the language picker and Home Assistant address and token fields.",
    light: "A ceiling light page with a brightness slider.",
    cover:
      "A curtain page with a position slider and open, stop and close buttons.",
    media: "A TV page with pause, previous and next track buttons.",
  },
  ui: {
    skip: "Skip to content",
    language: "Language",
    featureVisual: "Decorative illustration of the NeoHM interface",
    menuBar: "NeoHM in the macOS menu bar",
    primaryNav: "Primary",
  },
}
