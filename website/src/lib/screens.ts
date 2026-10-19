import coverEn from "@/assets/screens/en/cover-curtain.webp"
import homeEn from "@/assets/screens/en/home.webp"
import lightEn from "@/assets/screens/en/light-ceiling.webp"
import mediaEn from "@/assets/screens/en/media_player-tv.webp"
import climateEn from "@/assets/screens/en/panel.webp"
import settingsEn from "@/assets/screens/en/settings.webp"
import coverJa from "@/assets/screens/ja/cover-curtain.webp"
import homeJa from "@/assets/screens/ja/home.webp"
import lightJa from "@/assets/screens/ja/light-ceiling.webp"
import mediaJa from "@/assets/screens/ja/media_player-tv.webp"
import climateJa from "@/assets/screens/ja/panel.webp"
import settingsJa from "@/assets/screens/ja/settings.webp"
import coverKo from "@/assets/screens/ko/cover-curtain.webp"
import homeKo from "@/assets/screens/ko/home.webp"
import lightKo from "@/assets/screens/ko/light-ceiling.webp"
import mediaKo from "@/assets/screens/ko/media_player-tv.webp"
import climateKo from "@/assets/screens/ko/panel.webp"
import settingsKo from "@/assets/screens/ko/settings.webp"
import coverZhHans from "@/assets/screens/zh-Hans/cover-curtain.webp"
import homeZhHans from "@/assets/screens/zh-Hans/home.webp"
import lightZhHans from "@/assets/screens/zh-Hans/light-ceiling.webp"
import mediaZhHans from "@/assets/screens/zh-Hans/media_player-tv.webp"
import climateZhHans from "@/assets/screens/zh-Hans/panel.webp"
import settingsZhHans from "@/assets/screens/zh-Hans/settings.webp"
import coverZhHant from "@/assets/screens/zh-Hant/cover-curtain.webp"
import homeZhHant from "@/assets/screens/zh-Hant/home.webp"
import lightZhHant from "@/assets/screens/zh-Hant/light-ceiling.webp"
import mediaZhHant from "@/assets/screens/zh-Hant/media_player-tv.webp"
import climateZhHant from "@/assets/screens/zh-Hant/panel.webp"
import settingsZhHant from "@/assets/screens/zh-Hant/settings.webp"

export const shots = {
  en: {
    home: homeEn,
    climate: climateEn,
    settings: settingsEn,
    light: lightEn,
    cover: coverEn,
    media: mediaEn,
  },
  "zh-Hans": {
    home: homeZhHans,
    climate: climateZhHans,
    settings: settingsZhHans,
    light: lightZhHans,
    cover: coverZhHans,
    media: mediaZhHans,
  },
  "zh-Hant": {
    home: homeZhHant,
    climate: climateZhHant,
    settings: settingsZhHant,
    light: lightZhHant,
    cover: coverZhHant,
    media: mediaZhHant,
  },
  ja: {
    home: homeJa,
    climate: climateJa,
    settings: settingsJa,
    light: lightJa,
    cover: coverJa,
    media: mediaJa,
  },
  ko: {
    home: homeKo,
    climate: climateKo,
    settings: settingsKo,
    light: lightKo,
    cover: coverKo,
    media: mediaKo,
  },
} as const

export type ShotName = keyof (typeof shots)["en"]
