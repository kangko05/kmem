export const StorageBar = ({
  percentage,
  used,
  total,
}: {
  percentage: number;
  used: string;
  total: string;
}) => {
  const getGradient = (pct: number) => {
    if (pct < 70) return "bg-green-500";
    if (pct < 85) return "bg-yellow-500";
    if (pct < 95) return "bg-orange-500";
    return "bg-red-500";
  };

  return (
    <div className="w-full sm:max-w-xl storage-usage">
      <div className="flex justify-between mb-2">
        <span className="text-sm font-medium">storage usage</span>
        <span className="text-sm text-gray-500">
          {used} / {total}
        </span>
      </div>
      <div className="w-full bg-gray-200 rounded-lg h-3 relative overflow-hidden">
        <div
          className={`h-3 rounded-lg transition-all duration-500 ${getGradient(percentage)} absolute top-0 left-0`}
          style={{ width: `${Math.max(percentage, 0)}%` }}
        />
      </div>
      <div className="flex justify-between text-xs text-gray-500 my-1">
        <span>{Math.floor(percentage)}% used</span>
      </div>
    </div>
  );
};
