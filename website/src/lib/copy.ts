import { en } from "@/lib/copy/en"
import { ja } from "@/lib/copy/ja"
import { ko } from "@/lib/copy/ko"
import { zhHans } from "@/lib/copy/zh-Hans"
import { zhHant } from "@/lib/copy/zh-Hant"

export type Copy = typeof en

export const copy = {
  en,
  "zh-Hans": zhHans,
  "zh-Hant": zhHant,
  ja,
  ko,
} satisfies Record<string, Copy>

export type Language = keyof typeof copy

export const languageOptions = [
  { value: "en", label: "English" },
  { value: "zh-Hans", label: "简体中文" },
  { value: "zh-Hant", label: "繁體中文" },
  { value: "ja", label: "日本語" },
  { value: "ko", label: "한국어" },
] as const satisfies ReadonlyArray<{ value: Language; label: string }>
