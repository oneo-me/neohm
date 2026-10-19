import appIcon from "@/assets/app-icon.webp"

export function Logo() {
  return (
    <span className="flex items-center gap-2.5">
      <img
        src={appIcon}
        alt=""
        width={512}
        height={512}
        className="size-8 rounded-[9px] shadow-[0_4px_16px_rgba(0,0,0,0.35)]"
      />
      <span className="text-[15px] font-semibold tracking-tight">NeoHM</span>
    </span>
  )
}
