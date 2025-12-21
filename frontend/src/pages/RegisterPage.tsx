import { Link } from "react-router-dom";
import { RegisterForm } from "../features/auth/components/RegisterForm";

export const RegisterPage = () => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-slate-950 p-4">
      <RegisterForm />
      <p className="mt-6 text-slate-400 text-sm">
        ¿Ya tienes cuenta?{" "}
        <Link to="/login" className="text-emerald-500 hover:underline">
          Inicia sesión
        </Link>
      </p>
    </div>
  );
};
