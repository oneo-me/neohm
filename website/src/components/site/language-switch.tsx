import { Globe2 } from "lucide-react"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { languageOptions } from "@/lib/copy"
import { useLanguage } from "@/lib/language"
import { cn } from "@/lib/utils"

export function LanguageSwitch({ className }: { className?: string }) {
  const { lang, setLang, t } = useLanguage()

  return (
    <Select
      value={lang}
      onValueChange={(value) => {
        if (value) {
          setLang(value)
        }
      }}
      items={languageOptions}
    >
      <SelectTrigger
        aria-label={t.ui.language}
        className={cn(
          "h-8 w-auto gap-1.5 rounded-full border-white/10 bg-white/[0.04] py-0 pr-2 pl-2.5 text-[11px] font-medium text-muted-foreground transition-colors hover:bg-white/[0.06] hover:text-foreground focus-visible:ring-ring/40 data-[size=sm]:h-8",
          className
        )}
      >
        <Globe2 className="size-3.5" />
        <SelectValue />
      </SelectTrigger>
      <SelectContent
        align="end"
        sideOffset={8}
        alignItemWithTrigger={false}
        className="min-w-44 rounded-xl border border-white/10 bg-popover/95 p-1 shadow-2xl backdrop-blur-xl"
      >
        {languageOptions.map((language) => (
          <SelectItem
            key={language.value}
            value={language.value}
            className="rounded-lg py-1.5 pr-8 pl-2.5 text-[13px] focus:bg-white/10 focus:text-foreground"
          >
            {language.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
