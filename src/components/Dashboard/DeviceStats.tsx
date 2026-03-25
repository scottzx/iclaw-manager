import { useState, useEffect, useCallback } from 'react';
import {
  Cpu,
  MemoryStick,
  HardDrive,
  Activity,
  RefreshCw,
  AlertCircle,
  CheckCircle2,
  XCircle,
} from 'lucide-react';
import { api } from '../../lib/tauri';
import type { SystemInfo, SystemUsage, SystemStatus } from '../../lib/types';

function StatCard({
  icon: Icon,
  label,
  value,
  unit,
  colorClass,
}: {
  icon: React.ElementType;
  label: string;
  value: string | number;
  unit?: string;
  colorClass: string;
}) {
  return (
    <div className="bg-surface-card rounded-2xl p-5 border border-edge">
      <div className="flex items-start gap-4">
        <div
          className={`w-10 h-10 rounded-xl flex items-center justify-center shrink-0 ${colorClass}`}
        >
          <Icon size={20} className="text-white" />
        </div>
        <div>
          <div className="text-2xl font-bold tabular-nums text-content-primary">
            {value}
            {unit && (
              <span className="text-sm font-normal text-content-tertiary">{unit}</span>
            )}
          </div>
          <div className="text-sm mt-0.5 text-content-secondary">{label}</div>
        </div>
      </div>
    </div>
  );
}

function ServiceRow({
  name,
  active,
  running,
}: {
  name: string;
  active: boolean;
  running: boolean;
}) {
  return (
    <div
      className="flex items-center justify-between py-2 border-b border-edge"
      style={{ borderColor: 'var(--border-secondary)' }}
    >
      <div className="flex items-center gap-2">
        {running ? (
          <CheckCircle2 size={15} className="text-green-400" />
        ) : active ? (
          <AlertCircle size={15} className="text-amber-400" />
        ) : (
          <XCircle size={15} className="text-content-tertiary" />
        )}
        <span className="text-sm font-medium text-content-primary">{name}</span>
      </div>
      <span
        className="text-xs px-2 py-0.5 rounded-full"
        style={{
          backgroundColor: running
            ? 'rgba(74, 222, 128, 0.1)'
            : active
            ? 'rgba(251, 191, 36, 0.1)'
            : 'var(--bg-elevated)',
          color: running ? '#4ade80' : active ? '#fbbf24' : 'var(--text-tertiary)',
        }}
      >
        {running ? '运行中' : active ? '激活' : '停止'}
      </span>
    </div>
  );
}

function formatUptime(seconds: number): string {
  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  if (d > 0) return `${d}天 ${h}小时`;
  if (h > 0) return `${h}小时 ${m}分钟`;
  return `${m}分钟`;
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
}

interface DeviceStatsProps {
  compact?: boolean;
}

export function DeviceStats({ compact = false }: DeviceStatsProps) {
  const [info, setInfo] = useState<SystemInfo | null>(null);
  const [usage, setUsage] = useState<SystemUsage | null>(null);
  const [status, setStatus] = useState<SystemStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      setError(null);
      const [infoData, usageData, statusData] = await Promise.all([
        api.getDeviceSystemInfo() as Promise<SystemInfo>,
        api.getSystemUsage() as Promise<SystemUsage>,
        api.getSystemStatus() as Promise<SystemStatus>,
      ]);
      setInfo(infoData);
      setUsage(usageData);
      setStatus(statusData);
    } catch (err) {
      setError('无法获取设备状态');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, [fetchData]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-32">
        <RefreshCw size={20} className="animate-spin text-content-tertiary" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-32 text-content-secondary">
        <AlertCircle size={20} className="mr-2" />
        <span>{error}</span>
      </div>
    );
  }

  if (compact) {
    // 紧凑模式：只显示4个小卡片
    return (
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <StatCard
          icon={Cpu}
          label="CPU"
          value={usage?.cpu_percent?.toFixed(0) ?? 0}
          unit="%"
          colorClass="bg-indigo-500"
        />
        <StatCard
          icon={MemoryStick}
          label="内存"
          value={usage?.memory_percent?.toFixed(0) ?? 0}
          unit="%"
          colorClass="bg-purple-500"
        />
        <StatCard
          icon={HardDrive}
          label="磁盘"
          value={usage?.disk_percent?.toFixed(0) ?? 0}
          unit="%"
          colorClass="bg-amber-500"
        />
        <StatCard
          icon={Activity}
          label="运行时"
          value={info ? formatUptime(info.uptime) : '--'}
          colorClass="bg-emerald-500"
        />
      </div>
    );
  }

  return (
    <div className="space-y-6">

      {/* 详情行 */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* 资源进度条 */}
        <div className="lg:col-span-2 bg-surface-card rounded-2xl p-5 border border-edge space-y-4">
          <h3 className="text-sm font-semibold text-content-primary">资源使用详情</h3>

          <div>
            <div className="flex justify-between text-sm mb-1.5">
              <span className="text-content-secondary">CPU</span>
              <span className="font-medium text-content-primary">
                {usage?.cpu_percent?.toFixed(1) ?? 0}%
              </span>
            </div>
            <div className="h-2 rounded-full overflow-hidden bg-surface-elevated">
              <div
                className="h-full bg-indigo-500 rounded-full transition-all duration-500"
                style={{ width: `${usage?.cpu_percent ?? 0}%` }}
              />
            </div>
            <div className="text-xs mt-1 text-content-tertiary">{info?.cpu_model}</div>
          </div>

          <div>
            <div className="flex justify-between text-sm mb-1.5">
              <span className="text-content-secondary">内存</span>
              <span className="font-medium text-content-primary">
                {usage
                  ? `${formatBytes(usage.memory_used)} / ${formatBytes(usage.memory_total)}`
                  : '--'}
              </span>
            </div>
            <div className="h-2 rounded-full overflow-hidden bg-surface-elevated">
              <div
                className="h-full bg-purple-500 rounded-full transition-all duration-500"
                style={{ width: `${usage?.memory_percent ?? 0}%` }}
              />
            </div>
          </div>

          <div>
            <div className="flex justify-between text-sm mb-1.5">
              <span className="text-content-secondary">磁盘</span>
              <span className="font-medium text-content-primary">
                {usage
                  ? `${(usage.disk_used / 1024 ** 3).toFixed(1)} GB / ${(usage.disk_total / 1024 ** 3).toFixed(1)} GB`
                  : '--'}
              </span>
            </div>
            <div className="h-2 rounded-full overflow-hidden bg-surface-elevated">
              <div
                className="h-full bg-amber-500 rounded-full transition-all duration-500"
                style={{ width: `${usage?.disk_percent ?? 0}%` }}
              />
            </div>
          </div>
        </div>

        {/* 服务状态列表 */}
        <div className="bg-surface-card rounded-2xl p-5 border border-edge">
          <h3 className="text-sm font-semibold mb-3 text-content-primary">服务状态</h3>
          <div className="divide-y" style={{ borderColor: 'var(--border-secondary)' }}>
            {status?.services?.map((svc) => (
              <ServiceRow
                key={svc.name}
                name={
                  svc.name === 'openclaw'
                    ? 'OpenClaw AI'
                    : svc.name === 'adminApi'
                    ? 'Admin API'
                    : svc.name
                }
                active={svc.active}
                running={svc.running}
              />
            ))}
          </div>
        </div>
      </div>

    </div>
  );
}
