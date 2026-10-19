export const zhHant = {
  meta: {
    title: "NeoHM — 選單列裡的 Home Assistant",
    description:
      "NeoHM 是一款選單列與系統匣應用程式，把 Home Assistant 的裝置依房間排成卡片，從選單列或系統匣開啟面板控制空調，為 macOS、Windows 和 Linux 提供原生版本。",
  },
  nav: {
    features: "功能",
    screens: "細看面板",
    setup: "安裝",
    faq: "常見問題",
    github: "GitHub",
    download: "下載",
    language: "語言",
  },
  hero: {
    badge: "為 Home Assistant 打造的選單列應用程式",
    titleLead: "你的家，",
    titleAccent: "就在選單列裡。",
    body: "NeoHM 把 Home Assistant 的每個房間都放到一次點擊之外。開關燈、拉窗簾、調空調，不必離開手邊的事。",
    primary: "下載 NeoHM",
    secondary: "在 GitHub 查看",
    meta: ["macOS · Windows · Linux", "全架構版本", "5 種語言"],
    note: "免費開源，支援 macOS、Windows 和 Linux。",
  },
  features: {
    eyebrow: "為什麼選 NeoHM",
    title: "一眼看清，不打斷手邊的事。",
    body: "面板呈現家裡的狀態，控制則深入到每台裝置真正支援的功能。",
    items: [
      {
        title: "房間優先",
        body: "裝置以卡片形式排在所屬房間下，開啟時帶著柔和的光。分組來自 Home Assistant 的區域與裝置登錄表。",
      },
      {
        title: "點一下就能開關",
        body: "點擊卡片開關裝置；不能開關的裝置會直接開啟頁面。角落的箭頭總能更進一步：調光、調色、移動、播放或上鎖。",
      },
      {
        title: "完整的空調控制",
        body: "電源、溫度、模式、風速、風向、節能模式和睡眠，裝置支援多少就顯示多少。",
      },
      {
        title: "安靜的選單列",
        body: "無論裝置是否開啟，圖示始終如一，旁邊也不顯示其他資訊。需要時點一下開啟面板，再點一下收起。",
      },
      {
        title: "設定也在面板裡",
        body: "語言、位址和權杖都是面板裡的一頁，不會另開視窗。可選 English、简体中文、繁體中文、日本語、한국어。",
      },
      {
        title: "每種裝置都有頁面",
        body: "燈、插座、風扇、窗簾、門鎖、媒體播放器和感測器，各自擁有真正用得到的檔位、顏色、操作和讀數。",
      },
    ],
  },
  showcase: {
    eyebrow: "細看面板",
    title: "小巧的面板，完整的控制。",
    body: "面板按你家的樣子展開。下面是應用程式的真實螢幕截圖，不是示意圖。",
    tabs: {
      home: {
        label: "首頁",
        title: "依房間看整個家",
        body: "每台裝置都是所屬房間下的一張卡片。亮起的卡片正在運作，名稱下面一行是真正重要的狀態。",
        points: [
          "依 Home Assistant 的區域分組",
          "狀態、亮度、位置或溫度",
          "點擊箭頭進入裝置頁面",
        ],
      },
      climate: {
        label: "空調",
        title: "一個真正的空調遙控器",
        body: "室內溫度在設定溫度上方。下方是滑桿和裝置真正提供的控制，不多也不少。",
        points: [
          "目標溫度以 1° 為步長",
          "模式、風速、風向、節能與睡眠",
          "名稱下方顯示室內溫度",
        ],
      },
      devices: {
        label: "裝置",
        title: "每台裝置都有自己的頁面",
        body: "燈調節亮度和顏色，窗簾調節位置，電視控制播放，門鎖控制開關。頁面的樣子由裝置決定。",
        points: [
          "亮度、顏色、位置、風速或音量",
          "開啟、關閉、上鎖、解鎖、播放或暫停",
          "色板選擇顏色，或用取色器調出任意顏色",
          "窗簾或閥門可直接跳到 0、25、50、75 或 100%",
        ],
        shots: ["燈", "窗簾", "電視"],
      },
      settings: {
        label: "設定",
        title: "連線一次，然後忘掉它",
        body: "一個位址、一個權杖、一種語言。設定是面板裡的一頁，點齒輪開啟，不會變成另一個視窗。",
        points: [
          "長期存取權杖與你的位址",
          "預設跟隨系統的語言",
          "一切都留在你的電腦上",
        ],
      },
    },
  },
  languages: {
    eyebrow: "多語言",
    title: "五種語言，同一個面板。",
    body: "NeoHM 預設跟隨系統的語言，也可以在設定裡自行選擇。",
  },
  setup: {
    eyebrow: "安裝",
    title: "從下載到選單列。",
    body: "NeoHM 支援 macOS 13 或以上、Windows 和 Linux。下載安裝，建立存取權杖，你的家就在一次點擊之外。",
    steps: [
      {
        title: "下載 NeoHM",
        body: "在 GitHub 取得最新版本。macOS 上把 NeoHM 拖進「應用程式」資料夾；它沒有經過 Apple 公證，首次啟動前請在終端機執行 xattr -dr com.apple.quarantine /Applications/NeoHM.app。Windows 上執行安裝程式。Linux 上執行 curl -fsSL https://github.com/oneo-me/neohm/releases/latest/download/install.sh | sh。",
      },
      {
        title: "建立存取權杖",
        body: "在 Home Assistant 的個人資料中建立長期存取權杖。",
      },
      {
        title: "連線，關閉設定",
        body: "開啟 NeoHM，填入位址和權杖並連線，房間就在眼前。",
      },
    ],
    downloadTitle: "NeoHM macOS、Windows 與 Linux 版",
    downloadBody: "最新版本 · macOS 13 或以上、Windows、Linux",
    download: "從 GitHub 下載",
    releases: "查看所有版本",
    free: "免費開源",
    platforms: [
      { name: "macOS", status: ".dmg" },
      { name: "Windows", status: "安裝程式" },
      { name: "Linux", status: "install.sh · .deb" },
    ],
    architectures: "Apple 晶片 · Intel · x64 · ARM64",
    privacyTitle: "資料留在本機",
    privacyBody:
      "位址、權杖和語言儲存在 ~/.neohm/config.json，只有你的使用者能讀取。NeoHM 只與你填寫的 Home Assistant 通訊，僅在檢查更新時存取 GitHub，預設每天一次，可在設定中關閉。",
    requirements: ["macOS 13+", "Windows · Linux", "全架構", "Home Assistant"],
  },
  faq: {
    eyebrow: "答疑",
    title: "常見問題",
    body: "簡短的答案：一個 Home Assistant，一台電腦，中間沒有別的東西。",
    more: "還有別的疑問？",
    link: "在 GitHub 上提問",
    items: [
      {
        question: "NeoHM 免費嗎？",
        answer: "免費，而且開源。在 GitHub 下載最新版本，原始碼也在那裡。",
      },
      {
        question: "支援哪些平台？",
        answer:
          "macOS 13 或以上、Windows 和 Linux，並提供全架構版本：macOS 的 Apple 晶片與 Intel，Windows 和 Linux 的 x64 與 ARM64。Linux 上的系統匣圖示需要 libayatana-appindicator3。",
      },
      {
        question: "需要怎樣的 Home Assistant 環境？",
        answer:
          "只要電腦能連線即可，區域網路或遠端位址都可以。NeoHM 透過 REST 讀取裝置狀態，透過 Home Assistant 的 WebSocket API 讀取區域與裝置登錄表。",
      },
      {
        question: "存取權杖儲存在哪裡？",
        answer:
          "和其他設定一起儲存在 ~/.neohm/config.json，只有你的使用者能讀取。權杖只會傳送到你填寫的位址。",
      },
      {
        question: "可以連接多個家嗎？",
        answer:
          "不行，NeoHM 一次控制一個 Home Assistant 執行個體。設定裡的位址和權杖決定是哪一個。",
      },
      {
        question: "多久重新整理一次？",
        answer:
          "面板開啟時每 3 秒重新整理一次；面板收起時不會讀取，你的操作仍會立即顯示。",
      },
      {
        question: "選單列圖示建立失敗怎麼辦？",
        answer: "NeoHM 會在一般視窗中開啟同樣的頁面，功能不受影響。",
      },
      {
        question: "支援哪些語言？",
        answer: "English、简体中文、繁體中文、日本語、한국어，預設跟隨系統。",
      },
    ],
  },
  cta: {
    title: ["開啟面板，", "不必切換應用程式。"],
    body: "在 NeoHM 裡的每一次變更都會直接傳送到 Home Assistant，沒有其他需要同步的東西。",
    primary: "下載 NeoHM",
    secondary: "閱讀原始碼",
  },
  footer: {
    tagline: "屬於 Home Assistant 的選單列。",
    note: "NeoHM 是獨立專案，與 Home Assistant 或 Nabu Casa 沒有關聯。",
    links: {
      github: "GitHub",
      readme: "README",
      architecture: "架構",
      changelog: "更新日誌",
    },
    project: "專案",
    builtWith: "使用 Go、MyGo 與 Onevra 打造。",
    madeBy: "為 macOS、Windows 與 Linux 而做。",
  },
  alts: {
    home: "NeoHM 面板，Home Assistant 裝置依房間分組，吸頂燈和窗簾正在開啟。",
    climate: "NeoHM 的空調頁面，包含溫度滑桿、節能模式和模式控制。",
    settings: "NeoHM 設定，包含語言選擇以及 Home Assistant 位址和權杖輸入框。",
    light: "吸頂燈頁面，包含亮度滑桿。",
    cover: "窗簾頁面，包含位置滑桿以及開啟、停止和關閉按鈕。",
    media: "電視頁面，包含暫停、上一首和下一首按鈕。",
  },
  ui: {
    skip: "跳到主要內容",
    language: "語言",
    featureVisual: "NeoHM 介面裝飾插圖",
    menuBar: "macOS 選單列中的 NeoHM",
    primaryNav: "主要導覽",
  },
}
