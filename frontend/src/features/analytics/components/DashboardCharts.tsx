import { 
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell, ReferenceLine
} from 'recharts';

interface SportPerformance {
    sport_key: string;
    bets: number;
    profit: number;
}

interface Props {
    data: SportPerformance[];
}

export const ProfitBySportChart = ({ data }: Props) => {
    // Si no hay datos, mostramos mensaje
    if (!data || data.length === 0) {
        return <div className="h-64 flex items-center justify-center text-slate-500">Sin datos suficientes aún</div>;
    }

    return (
        <div className="bg-slate-900 p-6 rounded-xl border border-slate-800 shadow-lg">
            <h3 className="text-lg font-bold text-white mb-4">Rendimiento por Deporte ($)</h3>
            <div className="h-64 w-full">
                <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data}>
                        <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                        <XAxis dataKey="sport_key" stroke="#94a3b8" />
                        <YAxis stroke="#94a3b8" />
                        <Tooltip 
                            contentStyle={{ backgroundColor: '#0f172a', borderColor: '#334155', color: '#fff' }}
                            formatter={(value: number) => [`$${value.toFixed(2)}`, 'Ganancia/Pérdida']}
                        />
                        <ReferenceLine y={0} stroke="#475569" />
                        <Bar dataKey="profit" fill="#10b981" radius={[4, 4, 0, 0]}>
                            {data.map((entry, index) => (
                                <Cell key={`cell-${index}`} fill={entry.profit >= 0 ? '#10b981' : '#ef4444'} />
                            ))}
                        </Bar>
                    </BarChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
};