export const ko = {
  meta: {
    title: "NeoHM — 메뉴 막대의 Home Assistant",
    description:
      "Home Assistant를 위한 메뉴 막대 앱입니다. 기기를 방별로 전환하고, 메뉴 막대에서 여는 패널로 에어컨을 제어하세요. macOS는 모든 아키텍처를 지원하며, Windows와 Linux도 모든 아키텍처를 지원할 예정입니다.",
  },
  nav: {
    features: "기능",
    screens: "자세히 보기",
    setup: "설치",
    faq: "FAQ",
    github: "GitHub",
    download: "다운로드",
    language: "언어",
  },
  hero: {
    badge: "Home Assistant를 위한 메뉴 막대 앱",
    titleLead: "우리 집을,",
    titleAccent: "메뉴 막대에.",
    body: "NeoHM은 Home Assistant의 모든 방을 클릭 한 번 거리에 둡니다. 조명을 켜고, 커튼을 열고, 에어컨을 조절하세요. 하던 일은 그대로 이어집니다.",
    primary: "NeoHM 다운로드",
    secondary: "GitHub에서 보기",
    meta: ["macOS · Windows · Linux", "모든 아키텍처", "5개 언어"],
    note: "무료 오픈 소스. macOS 버전은 지금, Windows와 Linux 버전은 곧 만나보세요.",
  },
  features: {
    eyebrow: "NeoHM을 선택하는 이유",
    title: "한눈에 보이고, 방해하지 않습니다.",
    body: "집의 상태를 보여 주는 패널과, 각 기기가 실제로 지원하는 범위까지 이어지는 제어.",
    items: [
      {
        title: "방이 먼저",
        body: "기기는 해당 구역 아래 타일로 표시되고, 켜져 있으면 따뜻한 빛을 띱니다. 그룹은 Home Assistant 레지스트리에서 가져옵니다.",
      },
      {
        title: "클릭해서 전환",
        body: "타일을 클릭하면 기기를 켜고 끕니다. 스위치가 없는 기기는 페이지를 엽니다. 모서리의 화살표는 밝기, 색상, 이동, 재생, 잠금으로 이어집니다.",
      },
      {
        title: "에어컨은 전부",
        body: "전원, 온도, 모드, 풍량, 풍향, 절전 모드, 취침. 기기가 지원하는 만큼만 보여 줍니다.",
      },
      {
        title: "조용한 메뉴 막대",
        body: "기기 상태와 상관없이 아이콘은 늘 같고, 옆에는 아무것도 표시하지 않습니다. 필요할 때 클릭해 패널을 열고, 다시 클릭하면 닫힙니다.",
      },
      {
        title: "설정도 패널 안에서",
        body: "언어, 주소, 토큰은 패널의 한 페이지이며 별도 창을 열지 않습니다. English, 简体中文, 繁體中文, 日本語, 한국어 중에서 선택할 수 있습니다.",
      },
      {
        title: "모든 기기에 페이지를",
        body: "조명, 플러그, 선풍기, 커튼, 잠금 장치, 미디어 플레이어, 센서가 실제로 지원하는 단계, 색상, 동작, 측정값을 제공합니다.",
      },
    ],
  },
  showcase: {
    eyebrow: "자세히 보기",
    title: "작은 패널, 온전한 제어.",
    body: "패널은 우리 집의 모양을 따릅니다. 아래는 앱의 실제 스크린샷입니다.",
    tabs: {
      home: {
        label: "홈",
        title: "방별로 보는 집",
        body: "모든 기기가 해당 구역 아래 타일로 표시됩니다. 켜진 타일은 빛나고, 이름 아래 줄에 중요한 상태가 나타납니다.",
        points: [
          "Home Assistant 구역으로 그룹화",
          "상태, 밝기, 위치 또는 온도",
          "화살표로 기기 페이지 열기",
        ],
      },
      climate: {
        label: "에어컨",
        title: "진짜 같은 에어컨 리모컨",
        body: "설정 온도 위에 실내 온도가 표시됩니다. 그 아래에는 기기가 제공하는 슬라이더와 제어만 있습니다.",
        points: [
          "설정 온도는 1° 단위",
          "모드, 풍량, 풍향, 절전, 취침",
          "이름 아래 실내 온도",
        ],
      },
      devices: {
        label: "기기",
        title: "모든 기기에 전용 페이지",
        body: "조명은 밝기와 색상, 커튼은 위치, TV는 재생, 잠금 장치는 잠금을 제공합니다. 페이지의 모양은 기기가 정합니다.",
        points: [
          "밝기, 색상, 위치, 풍속 또는 음량",
          "열기, 닫기, 잠금, 해제, 재생 또는 일시 정지",
          "색상 견본 또는 피커에서 원하는 색 선택",
          "커튼과 밸브는 0, 25, 50, 75, 100%로 바로 이동",
        ],
        shots: ["조명", "커튼", "TV"],
      },
      settings: {
        label: "설정",
        title: "한 번 연결하고 잊으세요",
        body: "주소 하나, 토큰 하나, 언어 하나. 설정은 톱니바퀴로 여는 패널의 한 페이지이며, 별도 창이 되지 않습니다.",
        points: [
          "장기 액세스 토큰과 주소",
          "기본값은 시스템 언어를 따름",
          "모든 것이 컴퓨터 안에",
        ],
      },
    },
  },
  languages: {
    eyebrow: "언어",
    title: "다섯 개 언어, 하나의 패널.",
    body: "NeoHM은 기본적으로 컴퓨터의 언어를 따르며, 설정에서 선택할 수도 있습니다.",
  },
  setup: {
    eyebrow: "설치",
    title: "다운로드에서 메뉴 막대까지.",
    body: "NeoHM은 macOS 13 이상에서 사용할 수 있으며, Windows와 Linux 버전도 준비 중입니다. 설치하고 액세스 토큰을 만들면 우리 집이 바로 눈앞에 있습니다.",
    steps: [
      {
        title: "NeoHM 다운로드",
        body: "GitHub에서 최신 macOS 릴리스를 받아 NeoHM을 응용 프로그램 폴더로 드래그하세요. Apple 공증을 받지 않았으므로 처음 실행하기 전에 터미널에서 xattr -dr com.apple.quarantine /Applications/NeoHM.app을(를) 실행하세요.",
      },
      {
        title: "액세스 토큰 만들기",
        body: "Home Assistant 프로필에서 장기 액세스 토큰을 만드세요.",
      },
      {
        title: "연결하고 설정 닫기",
        body: "NeoHM을 열고 주소와 토큰을 입력해 연결하세요. 방이 기다리고 있습니다.",
      },
    ],
    downloadTitle: "macOS용 NeoHM",
    downloadBody: "최신 릴리스 · macOS 13 이상",
    download: ".dmg 다운로드",
    releases: "모든 릴리스",
    free: "무료 오픈 소스",
    platforms: [
      { name: "macOS", status: "지금 사용 가능" },
      { name: "Windows", status: "준비 중" },
      { name: "Linux", status: "준비 중" },
    ],
    architectures: "Apple silicon · Intel · x64 · ARM64",
    privacyTitle: "설계부터 로컬",
    privacyBody:
      "주소, 토큰, 언어는 ~/.neohm/config.json에 저장되며 사용자만 읽을 수 있습니다. NeoHM은 입력한 Home Assistant와 통신하며, 업데이트를 확인할 때만 GitHub에 접속합니다. 확인은 하루에 한 번이며 설정에서 끌 수 있습니다.",
    requirements: [
      "macOS 13+",
      "Windows & Linux 준비 중",
      "모든 아키텍처",
      "Home Assistant",
    ],
  },
  faq: {
    eyebrow: "FAQ",
    title: "자주 묻는 질문",
    body: "짧게 말하면, Home Assistant 하나, 컴퓨터 하나, 그 사이에는 아무것도 없습니다.",
    more: "더 궁금한 점이 있나요?",
    link: "GitHub에서 이슈 열기",
    items: [
      {
        question: "NeoHM은 무료인가요?",
        answer:
          "네. 무료 오픈 소스입니다. GitHub에서 최신 릴리스를 다운로드할 수 있고, 소스 코드도 같은 곳에 있습니다.",
      },
      {
        question: "어떤 플랫폼을 지원하나요?",
        answer:
          "NeoHM은 현재 macOS 13 이상에서 사용할 수 있습니다. Windows와 Linux 버전도 준비 중이며, 각 플랫폼의 모든 아키텍처용 빌드를 제공할 예정입니다. macOS는 Apple silicon과 Intel, Windows와 Linux는 x64와 ARM64를 지원합니다.",
      },
      {
        question: "어떤 Home Assistant 환경이 필요한가요?",
        answer:
          "컴퓨터에서 접근할 수 있는 인스턴스라면 로컬 네트워크든 원격 주소든 괜찮습니다. NeoHM은 REST로 기기 상태를 읽고, Home Assistant의 WebSocket API로 구역과 기기 레지스트리를 불러옵니다.",
      },
      {
        question: "액세스 토큰은 어디에 저장되나요?",
        answer:
          "다른 설정과 함께 ~/.neohm/config.json에 저장되며 사용자만 읽을 수 있습니다. 토큰은 입력한 주소로만 전송됩니다.",
      },
      {
        question: "여러 집을 연결할 수 있나요?",
        answer:
          "NeoHM은 한 번에 하나의 Home Assistant 인스턴스를 제어합니다. 설정의 주소와 토큰이 어느 집인지 결정합니다.",
      },
      {
        question: "얼마나 자주 업데이트되나요?",
        answer:
          "패널이 열려 있을 때 3초마다 갱신합니다. 패널이 닫혀 있을 때는 읽지 않으며, 직접 한 변경은 즉시 표시됩니다.",
      },
      {
        question: "메뉴 막대 아이콘을 만들 수 없으면 어떻게 하나요?",
        answer:
          "같은 페이지를 일반 창에서 열기 때문에 그대로 사용할 수 있습니다.",
      },
      {
        question: "어떤 언어를 지원하나요?",
        answer:
          "English, 简体中文, 繁體中文, 日本語, 한국어를 지원합니다. 기본적으로 시스템 언어를 따릅니다.",
      },
    ],
  },
  cta: {
    title: ["패널을 여세요.", "앱 전환은 그만."],
    body: "NeoHM의 모든 변경은 Home Assistant로 곧바로 전송됩니다. 동기화할 다른 것은 없습니다.",
    primary: "NeoHM 다운로드",
    secondary: "소스 보기",
  },
  footer: {
    tagline: "Home Assistant를 위한 메뉴 막대.",
    note: "NeoHM은 독립 프로젝트이며 Home Assistant 또는 Nabu Casa와 제휴하지 않습니다.",
    links: {
      github: "GitHub",
      readme: "README",
      architecture: "아키텍처",
      changelog: "변경 로그",
    },
    project: "프로젝트",
    builtWith: "Go, MyGo, Onevra로 만들었습니다.",
    madeBy: "macOS, Windows, Linux를 위해.",
  },
  alts: {
    home: "Home Assistant 기기를 방별로 묶어 보여 주는 NeoHM 패널. 천장 조명과 커튼이 켜져 있습니다.",
    climate: "온도 슬라이더와 절전 모드, 모드 제어가 있는 NeoHM 에어컨 페이지.",
    settings: "언어 선택과 Home Assistant 주소, 토큰 입력란이 있는 NeoHM 설정.",
    light: "밝기 슬라이더가 있는 천장 조명 페이지.",
    cover: "위치 슬라이더와 열기, 정지, 닫기 버튼이 있는 커튼 페이지.",
    media: "일시 정지, 이전, 다음 버튼이 있는 TV 페이지.",
  },
  ui: {
    skip: "본문으로 건너뛰기",
    language: "언어",
    featureVisual: "NeoHM 인터페이스 장식 일러스트",
    menuBar: "macOS 메뉴 막대에 표시된 NeoHM",
    primaryNav: "주요 탐색",
  },
}
