import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

type SectionHeadingProps = {
  eyebrow: string
  title: string
  body?: string
  align?: "start" | "center"
  className?: string
}

export function SectionHeading({
  eyebrow,
  title,
  body,
  align = "start",
  className,
}: SectionHeadingProps) {
  return (
    <div
      className={cn(
        "flex flex-col gap-4",
        align === "center" && "items-center text-center",
        className
      )}
    >
      <Badge
        variant="outline"
        className="h-6 gap-2 rounded-full border-white/10 bg-white/[0.03] px-3 text-[10px] font-medium tracking-[0.16em] text-muted-foreground uppercase"
      >
        <span className="size-1.5 rounded-full bg-brand-blue shadow-[0_0_8px_2px_rgba(125,178,255,0.45)]" />
        {eyebrow}
      </Badge>
      <h2 className="max-w-3xl text-3xl font-semibold tracking-[-0.03em] text-balance sm:text-4xl lg:text-[2.9rem] lg:leading-[1.08]">
        {title}
      </h2>
      {body ? (
        <p
          className={cn(
            "max-w-2xl text-base leading-relaxed text-muted-foreground",
            align === "center" && "mx-auto"
          )}
        >
          {body}
        </p>
      ) : null}
    </div>
  )
}
