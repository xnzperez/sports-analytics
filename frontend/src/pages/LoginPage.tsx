import { Link } from "react-router-dom";
import { LoginForm } from "../features/auth/components/LoginForm";

export const LoginPage = () => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-slate-950 p-4">
      <LoginForm />
      <p className="mt-6 text-slate-400 text-sm">
        ¿No tienes cuenta?{" "}
        <Link to="/register" className="text-emerald-500 hover:underline">
          Regístrate aquí
        </Link>
      </p>
    </div>
  );
};
