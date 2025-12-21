import { cn } from "../../lib/utils";

interface Props {
  status: string;
}

export const StatusBadge = ({ status }: Props) => {
  const styles: Record<string, string> = {
    pending: "bg-yellow-500/10 text-yellow-500 border-yellow-500/20",
    WON: "bg-emerald-500/10 text-emerald-500 border-emerald-500/20",
    LOST: "bg-red-500/10 text-red-500 border-red-500/20",
    VOID: "bg-slate-500/10 text-slate-400 border-slate-500/20",
  };

  const label: Record<string, string> = {
    pending: "Pendiente",
    WON: "Ganada",
    LOST: "Perdida",
    VOID: "Anulada",
  };

  return (
    <span
      className={cn(
        "px-2.5 py-0.5 rounded-full text-xs font-medium border",
        styles[status] || styles.VOID
      )}
    >
      {label[status] || status}
    </span>
  );
};
