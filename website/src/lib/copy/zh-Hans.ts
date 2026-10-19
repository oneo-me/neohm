export const zhHans = {
  meta: {
    title: "NeoHM — 菜单栏里的 Home Assistant",
    description:
      "NeoHM 是一款菜单栏与托盘应用，把 Home Assistant 的设备按房间排成卡片，从菜单栏或托盘打开面板控制空调，为 macOS、Windows 和 Linux 提供原生版本。",
  },
  nav: {
    features: "功能",
    screens: "细看面板",
    setup: "安装",
    faq: "常见问题",
    github: "GitHub",
    download: "下载",
    language: "语言",
  },
  hero: {
    badge: "为 Home Assistant 打造的菜单栏应用",
    titleLead: "你的家，",
    titleAccent: "就在菜单栏里。",
    body: "NeoHM 把 Home Assistant 的每个房间都放到一次点击之外。开关灯、拉窗帘、调空调，不必离开手头的事。",
    primary: "下载 NeoHM",
    secondary: "在 GitHub 查看",
    meta: ["macOS · Windows · Linux", "全架构版本", "5 种语言"],
    note: "免费开源，支持 macOS、Windows 和 Linux。",
  },
  features: {
    eyebrow: "为什么选 NeoHM",
    title: "一眼看清，不打断手头的事。",
    body: "面板呈现家里的状态，控制则深入到每台设备真正支持的功能。",
    items: [
      {
        title: "房间优先",
        body: "设备以卡片形式排在所属房间下，开启时带着柔和的光。分组来自 Home Assistant 的区域与设备注册表。",
      },
      {
        title: "点一下就能开关",
        body: "点击卡片开关设备；不能开关的设备会直接打开页面。角落的箭头总能更进一步：调光、调色、移动、播放或上锁。",
      },
      {
        title: "完整的空调控制",
        body: "电源、温度、模式、风速、风向、节能模式和睡眠，设备支持多少就显示多少。",
      },
      {
        title: "安静的菜单栏",
        body: "无论设备是否开启，图标始终如一，旁边也不显示别的信息。需要时点一下打开面板，再点一下收起。",
      },
      {
        title: "设置也在面板里",
        body: "语言、地址和令牌都是面板里的一页，不会另开窗口。可选 English、简体中文、繁體中文、日本語、한국어。",
      },
      {
        title: "每种设备都有页面",
        body: "灯、插座、风扇、窗帘、门锁、媒体播放器和传感器，各自拥有真正用得到的档位、颜色、操作和读数。",
      },
    ],
  },
  showcase: {
    eyebrow: "细看面板",
    title: "小巧的面板，完整的控制。",
    body: "面板按你家的样子展开。下面是应用的真实截图，不是示意图。",
    tabs: {
      home: {
        label: "首页",
        title: "按房间看整个家",
        body: "每台设备都是所属房间下的一张卡片。亮起的卡片正在工作，名字下面一行是真正重要的状态。",
        points: [
          "按 Home Assistant 的区域分组",
          "状态、亮度、位置或温度",
          "点击箭头进入设备页面",
        ],
      },
      climate: {
        label: "空调",
        title: "一个真正的空调遥控器",
        body: "室内温度在设定温度上方。下方是滑杆和设备真正提供的控制，不多也不少。",
        points: [
          "目标温度以 1° 为步长",
          "模式、风速、风向、节能与睡眠",
          "名字下方显示室内温度",
        ],
      },
      devices: {
        label: "设备",
        title: "每台设备都有自己的页面",
        body: "灯调节亮度和颜色，窗帘调节位置，电视控制播放，门锁控制开关。页面的样子由设备决定。",
        points: [
          "亮度、颜色、位置、风速或音量",
          "打开、关闭、上锁、解锁、播放或暂停",
          "色板选择颜色，或用取色器调出任意颜色",
          "窗帘或阀门可直接跳到 0、25、50、75 或 100%",
        ],
        shots: ["灯", "窗帘", "电视"],
      },
      settings: {
        label: "设置",
        title: "连接一次，然后忘掉它",
        body: "一个地址、一个令牌、一种语言。设置是面板里的一页，点齿轮打开，不会变成另一个窗口。",
        points: [
          "长期访问令牌与你的地址",
          "默认跟随系统的语言",
          "一切都留在你的电脑上",
        ],
      },
    },
  },
  languages: {
    eyebrow: "多语言",
    title: "五种语言，同一个面板。",
    body: "NeoHM 默认跟随系统的语言，也可以在设置里自行选择。",
  },
  setup: {
    eyebrow: "安装",
    title: "从下载到菜单栏。",
    body: "NeoHM 支持 macOS 13 或更高版本、Windows 和 Linux。下载安装，创建访问令牌，你的家就在一次点击之外。",
    steps: [
      {
        title: "下载 NeoHM",
        body: "在 GitHub 获取最新版本。macOS 上把 NeoHM 拖进「应用程序」文件夹；它没有经过 Apple 公证，首次启动前请在终端运行 xattr -dr com.apple.quarantine /Applications/NeoHM.app。Windows 上运行安装程序。Linux 上运行 curl -fsSL https://github.com/oneo-me/neohm/releases/latest/download/install.sh | sh。",
      },
      {
        title: "创建访问令牌",
        body: "在 Home Assistant 的个人资料中创建长期访问令牌。",
      },
      {
        title: "连接，关闭设置",
        body: "打开 NeoHM，填入地址和令牌并连接，房间就在眼前。",
      },
    ],
    downloadTitle: "NeoHM macOS、Windows 与 Linux 版",
    downloadBody: "最新版本 · macOS 13 或更高、Windows、Linux",
    download: "从 GitHub 下载",
    releases: "查看所有版本",
    free: "免费开源",
    platforms: [
      { name: "macOS", status: ".dmg" },
      { name: "Windows", status: "安装程序" },
      { name: "Linux", status: "install.sh · .deb" },
    ],
    architectures: "Apple 芯片 · Intel · x64 · ARM64",
    privacyTitle: "数据留在本机",
    privacyBody:
      "地址、令牌和语言保存在 ~/.neohm/config.json，只有你的用户能读取。NeoHM 只与你填写的 Home Assistant 通信，仅在检查更新时访问 GitHub，默认每天一次，可在设置中关闭。",
    requirements: ["macOS 13+", "Windows · Linux", "全架构", "Home Assistant"],
  },
  faq: {
    eyebrow: "答疑",
    title: "常见问题",
    body: "简短的答案：一个 Home Assistant，一台电脑，中间没有别的东西。",
    more: "还有别的疑问？",
    link: "在 GitHub 上提问",
    items: [
      {
        question: "NeoHM 免费吗？",
        answer: "免费，而且开源。在 GitHub 下载最新版本，源码也在那里。",
      },
      {
        question: "支持哪些平台？",
        answer:
          "macOS 13 或更高版本、Windows 和 Linux，并提供全架构版本：macOS 的 Apple 芯片与 Intel，Windows 和 Linux 的 x64 与 ARM64。Linux 上的托盘图标需要 libayatana-appindicator3。",
      },
      {
        question: "需要怎样的 Home Assistant 环境？",
        answer:
          "只要电脑能访问即可，局域网或远程地址都可以。NeoHM 通过 REST 读取设备状态，通过 Home Assistant 的 WebSocket API 读取区域与设备注册表。",
      },
      {
        question: "访问令牌保存在哪里？",
        answer:
          "和其他设置一起保存在 ~/.neohm/config.json，只有你的用户能读取。令牌只会发送到你填写的地址。",
      },
      {
        question: "可以连接多个家吗？",
        answer:
          "不行，NeoHM 一次控制一个 Home Assistant 实例。设置里的地址和令牌决定是哪一个。",
      },
      {
        question: "多久刷新一次？",
        answer:
          "面板打开时每 3 秒刷新一次；面板收起时不读取，你的操作仍会立即显示。",
      },
      {
        question: "菜单栏图标创建失败怎么办？",
        answer: "NeoHM 会在普通窗口中打开同样的页面，功能不受影响。",
      },
      {
        question: "支持哪些语言？",
        answer: "English、简体中文、繁體中文、日本語、한국어，默认跟随系统。",
      },
    ],
  },
  cta: {
    title: ["打开面板，", "不必切换应用。"],
    body: "在 NeoHM 里的每一次更改都会直接发送到 Home Assistant，没有别的需要同步。",
    primary: "下载 NeoHM",
    secondary: "阅读源码",
  },
  footer: {
    tagline: "属于 Home Assistant 的菜单栏。",
    note: "NeoHM 是独立项目，与 Home Assistant 或 Nabu Casa 没有关联。",
    links: {
      github: "GitHub",
      readme: "README",
      architecture: "架构",
      changelog: "更新日志",
    },
    project: "项目",
    builtWith: "使用 Go、MyGo 与 Onevra 构建。",
    madeBy: "为 macOS、Windows 与 Linux 而做。",
  },
  alts: {
    home: "NeoHM 面板，Home Assistant 设备按房间分组，吸顶灯和窗帘正在开启。",
    climate: "NeoHM 的空调页面，包含温度滑杆、节能模式和模式控制。",
    settings: "NeoHM 设置，包含语言选择以及 Home Assistant 地址和令牌输入框。",
    light: "吸顶灯页面，包含亮度滑杆。",
    cover: "窗帘页面，包含位置滑杆以及打开、停止和关闭按钮。",
    media: "电视页面，包含暂停、上一首和下一首按钮。",
  },
  ui: {
    skip: "跳到主要内容",
    language: "语言",
    featureVisual: "NeoHM 界面装饰插图",
    menuBar: "macOS 菜单栏中的 NeoHM",
    primaryNav: "主导航",
  },
}
