import { cn } from "../../../lib/utils";

interface Props {
  password?: string;
}

export const PasswordRequirements = ({ password = "" }: Props) => {
  // Definimos las reglas visuales (deben coincidir con tu Zod Schema)
  const requirements = [
    { id: 1, label: "Mínimo 8 caracteres", regex: /.{8,}/ },
    { id: 2, label: "Una letra mayúscula", regex: /[A-Z]/ },
    { id: 3, label: "Una letra minúscula", regex: /[a-z]/ },
    { id: 4, label: "Un número", regex: /[0-9]/ },
    { id: 5, label: "Un carácter especial (@$!%*?&)", regex: /[^A-Za-z0-9]/ },
  ];

  return (
    <div className="mt-2 space-y-1.5 bg-slate-800/50 p-3 rounded-md border border-slate-700/50">
      <p className="text-xs font-medium text-slate-400 mb-2">
        La contraseña debe contener:
      </p>
      <ul className="space-y-1">
        {requirements.map((req) => {
          const isMet = req.regex.test(password);
          return (
            <li
              key={req.id}
              className="flex items-center gap-2 text-xs transition-colors duration-300"
            >
              {/* Círculo indicador */}
              <div
                className={cn(
                  "w-4 h-4 rounded-full flex items-center justify-center border transition-all duration-300",
                  isMet
                    ? "bg-emerald-500/20 border-emerald-500 text-emerald-500"
                    : "bg-slate-700 border-slate-600 text-transparent"
                )}
              >
                {/* Icono Check */}
                <svg
                  className="w-2.5 h-2.5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  strokeWidth="3"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>

              {/* Texto */}
              <span
                className={cn(
                  isMet
                    ? "text-slate-200 line-through decoration-emerald-500/50"
                    : "text-slate-500"
                )}
              >
                {req.label}
              </span>
            </li>
          );
        })}
      </ul>
    </div>
  );
};
