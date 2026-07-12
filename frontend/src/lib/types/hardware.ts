export interface HardwareMetrics {
	// A metric is `null` when the backend couldn't resolve that sensor
	// (renamed/missing in the hardware tree), as opposed to `undefined`
	// which means the whole payload is absent (PC offline).
	cpu_temp_tctl: number | null;
	cpu_temp_ccd1: number | null;
	cpu_power: number | null;
	gpu_temp: number | null;
	gpu_memory_temp: number | null;
	gpu_power: number | null;
	pc_name: string;
	timestamp: number;
}
