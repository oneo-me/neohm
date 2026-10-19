export const ja = {
  meta: {
    title: "NeoHM — Home Assistant をメニューバーに",
    description:
      "Home Assistant のためのメニューバーアプリ。デバイスを部屋ごとに切り替え、メニューバーで開くパネルからエアコンを操作できます。macOS 版はすべてのアーキテクチャに対応し、Windows と Linux 版もすべてのアーキテクチャに対応予定です。",
  },
  nav: {
    features: "機能",
    screens: "パネルを見る",
    setup: "インストール",
    faq: "FAQ",
    github: "GitHub",
    download: "ダウンロード",
    language: "言語",
  },
  hero: {
    badge: "Home Assistant のためのメニューバーアプリ",
    titleLead: "あなたの家を、",
    titleAccent: "メニューバーに。",
    body: "NeoHM は Home Assistant のすべての部屋をワンクリックの先に置きます。照明をつけ、カーテンを開け、エアコンを調整。作業の流れを止めずに。",
    primary: "NeoHM をダウンロード",
    secondary: "GitHub で見る",
    meta: ["macOS · Windows · Linux", "すべてのアーキテクチャ", "5 言語"],
    note: "無料でオープンソース。macOS 版は提供中、Windows と Linux 版は準備中です。",
  },
  features: {
    eyebrow: "NeoHM の特長",
    title: "ひと目で把握。邪魔をしない。",
    body: "家の状態を映すパネルと、各デバイスが対応する範囲まで踏み込める操作。",
    items: [
      {
        title: "部屋ごとに",
        body: "デバイスはエリアごとのタイルで並び、オンなら温かな光をまといます。グループ分けは Home Assistant のレジストリから取得します。",
      },
      {
        title: "クリックで切り替え",
        body: "タイルをクリックしてデバイスのオン・オフを切り替えます。スイッチのないデバイスはページを開きます。角の矢印からは、調光・色の変更・移動・再生・施錠へ進めます。",
      },
      {
        title: "空調を余すところなく",
        body: "電源、温度、モード、風量、風向、エコモード、おやすみ。本体が対応する機能だけを表示します。",
      },
      {
        title: "静かなメニューバー",
        body: "デバイスの状態にかかわらずアイコンは変わらず、その横に何も表示しません。必要なときにクリックしてパネルを開き、もう一度クリックすると閉じます。",
      },
      {
        title: "設定もパネルの中で",
        body: "言語、アドレス、トークンはパネルの 1 ページ。別ウィンドウにはなりません。English、简体中文、繁體中文、日本語、한국어 から選べます。",
      },
      {
        title: "すべてのデバイスにページを",
        body: "照明、プラグ、ファン、カーテン、ロック、メディアプレーヤー、センサー。それぞれが実際に使えるレベル、色、操作、読み取り値を備えます。",
      },
    ],
  },
  showcase: {
    eyebrow: "パネルを見る",
    title: "小さなパネルに、すべての操作を。",
    body: "パネルは家の形に沿って広がります。以下はアプリの実際のスクリーンショットです。",
    tabs: {
      home: {
        label: "ホーム",
        title: "部屋ごとに家を見渡す",
        body: "すべてのデバイスがエリアごとのタイルに並びます。光っているタイルはオンで、名前の下の行に大切な状態が表示されます。",
        points: [
          "Home Assistant のエリアでグループ化",
          "状態、明るさ、位置、温度",
          "矢印からデバイスのページへ",
        ],
      },
      climate: {
        label: "空調",
        title: "本物の空調リモコン",
        body: "設定温度の上に室温。その下にスライダーと、本体が提供する操作だけが並びます。",
        points: [
          "設定温度は 1° 刻み",
          "モード、風量、風向、エコ、おやすみ",
          "名前の下に室温",
        ],
      },
      devices: {
        label: "デバイス",
        title: "どのデバイスにも専用ページ",
        body: "照明は明るさと色、カーテンは位置、テレビは再生、ロックは施錠。ページの形はデバイスが決めます。",
        points: [
          "明るさ、色、位置、風速、音量",
          "開く、閉じる、施錠、解錠、再生、一時停止",
          "スウォッチから色を選び、ピッカーで任意の色に",
          "カーテンやバルブは 0、25、50、75、100% へ一発で",
        ],
        shots: ["照明", "カーテン", "テレビ"],
      },
      settings: {
        label: "設定",
        title: "一度つないだら、あとは忘れて",
        body: "アドレス 1 つ、トークン 1 つ、言語 1 つ。設定はギアから開くパネルの 1 ページで、別ウィンドウにはなりません。",
        points: [
          "長期アクセストークンとアドレス",
          "デフォルトはシステムの言語に追従",
          "すべてコンピュータの中に",
        ],
      },
    },
  },
  languages: {
    eyebrow: "言語",
    title: "5 言語、1 つのパネル。",
    body: "NeoHM はお使いのコンピュータの言語に従います。設定で選ぶこともできます。",
  },
  setup: {
    eyebrow: "インストール",
    title: "ダウンロードからメニューバーへ。",
    body: "NeoHM は macOS 13 以降に対応しており、Windows と Linux 版も準備中です。インストールしてアクセストークンを作成すれば、家はすぐそこです。",
    steps: [
      {
        title: "NeoHM をダウンロード",
        body: "GitHub から最新の macOS 版を入手し、NeoHM を「アプリケーション」フォルダにドラッグします。Apple の公証を受けていないため、初回起動の前にターミナルで xattr -dr com.apple.quarantine /Applications/NeoHM.app を実行してください。",
      },
      {
        title: "アクセストークンを作成",
        body: "Home Assistant のプロフィールで長期アクセストークンを作成します。",
      },
      {
        title: "接続して設定を閉じる",
        body: "NeoHM を開き、アドレスとトークンを入力して接続。部屋が待っています。",
      },
    ],
    downloadTitle: "macOS 版 NeoHM",
    downloadBody: "最新リリース · macOS 13 以降",
    download: ".dmg をダウンロード",
    releases: "すべてのリリース",
    free: "無料・オープンソース",
    platforms: [
      { name: "macOS", status: "提供中" },
      { name: "Windows", status: "準備中" },
      { name: "Linux", status: "準備中" },
    ],
    architectures: "Apple silicon · Intel · x64 · ARM64",
    privacyTitle: "ローカルに保存",
    privacyBody:
      "アドレス、トークン、言語は ~/.neohm/config.json に保存され、あなたのユーザーだけが読み取れます。NeoHM は入力した Home Assistant と通信し、GitHub にはアップデートの確認でのみアクセスします。確認は 1 日に 1 回で、設定でオフにできます。",
    requirements: [
      "macOS 13+",
      "Windows・Linux 準備中",
      "全アーキテクチャ",
      "Home Assistant",
    ],
  },
  faq: {
    eyebrow: "FAQ",
    title: "よくある質問",
    body: "短く言えば、Home Assistant 1 台、コンピュータ 1 台、その間に何もありません。",
    more: "ほかに質問がありますか？",
    link: "GitHub で Issue を開く",
    items: [
      {
        question: "NeoHM は無料ですか？",
        answer:
          "はい。無料でオープンソースです。GitHub から最新リリースをダウンロードできます。ソースコードも同じ場所にあります。",
      },
      {
        question: "対応するプラットフォームは？",
        answer:
          "NeoHM は現在 macOS 13 以降に対応しています。Windows と Linux 版は準備中で、各プラットフォームのすべてのアーキテクチャ向けビルドを提供する予定です。macOS は Apple silicon と Intel、Windows と Linux は x64 と ARM64 に対応します。",
      },
      {
        question: "どのような Home Assistant 環境が必要ですか？",
        answer:
          "お使いのコンピュータから到達できるインスタンスなら、ローカルネットワークでもリモートアドレスでもかまいません。NeoHM は REST でデバイスの状態を読み取り、Home Assistant の WebSocket API でエリアとデバイスのレジストリを読み込みます。",
      },
      {
        question: "アクセストークンはどこに保存されますか？",
        answer:
          "ほかの設定とともに ~/.neohm/config.json に保存され、あなたのユーザーだけが読み取れます。トークンは入力したアドレスにのみ送信されます。",
      },
      {
        question: "複数の家を接続できますか？",
        answer:
          "NeoHM は一度に 1 つの Home Assistant インスタンスを操作します。設定のアドレスとトークンで切り替わります。",
      },
      {
        question: "どのくらいの頻度で更新されますか？",
        answer:
          "パネルを開いている間は 3 秒ごとです。パネルを閉じている間は読み取らず、操作はすぐに表示されます。",
      },
      {
        question: "メニューバーのアイコンを作成できない場合は？",
        answer: "同じページを通常のウィンドウで開くので、そのまま使えます。",
      },
      {
        question: "対応言語は？",
        answer:
          "English、简体中文、繁體中文、日本語、한국어 に対応しています。デフォルトではシステムの言語に従います。",
      },
    ],
  },
  cta: {
    title: ["パネルを開く。", "アプリの切り替えは不要。"],
    body: "NeoHM の変更はすべて Home Assistant に直接送られます。ほかに同期するものはありません。",
    primary: "NeoHM をダウンロード",
    secondary: "ソースを読む",
  },
  footer: {
    tagline: "Home Assistant のためのメニューバー。",
    note: "NeoHM は独立したプロジェクトで、Home Assistant や Nabu Casa とは提携していません。",
    links: {
      github: "GitHub",
      readme: "README",
      architecture: "アーキテクチャ",
      changelog: "変更履歴",
    },
    project: "プロジェクト",
    builtWith: "Go、MyGo、Onevra で構築。",
    madeBy: "macOS、Windows、Linux のために。",
  },
  alts: {
    home: "Home Assistant のデバイスを部屋ごとに並べた NeoHM パネル。シーリングライトとカーテンがオンです。",
    climate:
      "温度スライダー、エコモード、モード操作を備えた NeoHM のエアコンページ。",
    settings:
      "言語の選択と Home Assistant のアドレス、トークン入力欄がある NeoHM の設定。",
    light: "明るさスライダーがあるシーリングライトのページ。",
    cover: "位置スライダーと開く、停止、閉じるボタンがあるカーテンのページ。",
    media: "一時停止、前へ、次へボタンがあるテレビのページ。",
  },
  ui: {
    skip: "コンテンツへスキップ",
    language: "言語",
    featureVisual: "NeoHM インターフェースの装飾イラスト",
    menuBar: "macOS のメニューバーに表示された NeoHM",
    primaryNav: "メインナビゲーション",
  },
}
