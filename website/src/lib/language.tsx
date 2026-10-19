import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react"

import { copy, type Copy, type Language } from "@/lib/copy"
import { shots } from "@/lib/screens"

export type { Language }

type LanguageContextValue = {
  lang: Language
  setLang: (lang: Language) => void
  t: Copy
  screens: (typeof shots)[Language]
}

const STORAGE_KEY = "neohm.language"

const documentLanguages: Record<Language, string> = {
  en: "en",
  "zh-Hans": "zh-Hans",
  "zh-Hant": "zh-Hant",
  ja: "ja",
  ko: "ko",
}

const openGraphLocales: Record<Language, string> = {
  en: "en_US",
  "zh-Hans": "zh_CN",
  "zh-Hant": "zh_TW",
  ja: "ja_JP",
  ko: "ko_KR",
}

const LanguageContext = createContext<LanguageContextValue | null>(null)

function isLanguage(value: string | null): value is Language {
  return value !== null && Object.hasOwn(copy, value)
}

function matchLanguage(tag: string): Language | null {
  const locale = tag.toLowerCase()

  if (locale.startsWith("zh")) {
    return locale.includes("hant") || /-(tw|hk|mo)(-|$)/.test(locale)
      ? "zh-Hant"
      : "zh-Hans"
  }

  if (locale.startsWith("ja")) {
    return "ja"
  }

  if (locale.startsWith("ko")) {
    return "ko"
  }

  if (locale.startsWith("en")) {
    return "en"
  }

  return null
}

function detectLanguage(): Language {
  if (typeof window === "undefined") {
    return "en"
  }

  const stored = window.localStorage.getItem(STORAGE_KEY)
  if (stored === "zh") {
    return "zh-Hans"
  }
  if (isLanguage(stored)) {
    return stored
  }

  const preferred =
    navigator.languages && navigator.languages.length > 0
      ? navigator.languages
      : [navigator.language]

  for (const tag of preferred) {
    const language = matchLanguage(tag)
    if (language) {
      return language
    }
  }

  return "en"
}

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [lang, setLang] = useState<Language>(detectLanguage)

  const t = copy[lang]
  const screens = shots[lang]

  useEffect(() => {
    document.documentElement.lang = documentLanguages[lang]
    document.title = t.meta.title

    const description = document.querySelector<HTMLMetaElement>(
      'meta[name="description"]'
    )
    if (description) {
      description.content = t.meta.description
    }

    const ogTitle = document.querySelector<HTMLMetaElement>(
      'meta[property="og:title"]'
    )
    if (ogTitle) {
      ogTitle.content = t.meta.title
    }

    const ogDescription = document.querySelector<HTMLMetaElement>(
      'meta[property="og:description"]'
    )
    if (ogDescription) {
      ogDescription.content = t.meta.description
    }

    const ogLocale = document.querySelector<HTMLMetaElement>(
      'meta[property="og:locale"]'
    )
    if (ogLocale) {
      ogLocale.content = openGraphLocales[lang]
    }

    const ogImageAlt = document.querySelector<HTMLMetaElement>(
      'meta[property="og:image:alt"]'
    )
    if (ogImageAlt) {
      ogImageAlt.content = t.meta.title
    }

    const twitterTitle = document.querySelector<HTMLMetaElement>(
      'meta[name="twitter:title"]'
    )
    if (twitterTitle) {
      twitterTitle.content = t.meta.title
    }

    const twitterDescription = document.querySelector<HTMLMetaElement>(
      'meta[name="twitter:description"]'
    )
    if (twitterDescription) {
      twitterDescription.content = t.meta.description
    }

    window.localStorage.setItem(STORAGE_KEY, lang)
  }, [lang, t])

  const value = useMemo(
    () => ({ lang, setLang, t, screens }),
    [lang, t, screens]
  )

  return (
    <LanguageContext.Provider value={value}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  const context = useContext(LanguageContext)

  if (!context) {
    throw new Error("useLanguage must be used within a LanguageProvider")
  }

  return context
}
