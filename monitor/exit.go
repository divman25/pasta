package monitor

import (
	"log"
	"time"

	"github.com/jjh2kiss/netlinkconnector/cnproc"
	"github.com/jjh2kiss/pstat/config"
	"github.com/jjh2kiss/pstat/stats"
	"github.com/jjh2kiss/pstat/system/process"
)

func processExitEvent(event *cnproc.ProcEvent, process_table *process.ProcessTable, stats_table *stats.StatsTable, config *config.Config, logger *log.Logger) error {
	ev, err := event.ExitEvent()
	if err != nil {
		return err
	}
	stats_table.Update(int(ev.ProcessPid), stats.STAT_EXIT)
	process := process_table.GetOrDefault(int(ev.ProcessPid))
	duration := time.Now().Sub(process.TimeStamp)

	if config.Quiet == false {
		if config.Events&event.What != 0 {
			name := process.Cmdline.CombinedString(process.KernelThread, config.Shortname, config.Dirstrip)
			log.Printf("exit %5d  %5d    %5.3f %s\n",
				ev.ProcessPid,
				ev.ExitCode,
				duration.Seconds(),
				name,
			)
		}
	}

	//post action
	process_table.Delete(int(ev.ProcessPid))
	return nil
}
