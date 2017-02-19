package cluster

func (cluster *Cluster) testSwitchOverNoReadOnlyNoRplCheck(conf string, test string) bool {
	if cluster.initTestCluster(conf, test) == false {
		return false
	}
	cluster.conf.RplChecks = false
	cluster.conf.MaxDelay = 0
	cluster.LogPrintf("TEST : Master is %s", cluster.master.URL)
	cluster.conf.ReadOnly = false
	for _, s := range cluster.servers {
		_, err := s.Conn.Exec("set global read_only=0")
		if err != nil {
			cluster.LogPrintf("ERROR : %s", err.Error())
			cluster.closeTestCluster(conf, test)
		}
	}
	SaveMasterURL := cluster.master.URL
	switchoverChan <- true
	cluster.waitFailoverEnd()
	cluster.LogPrintf("TEST : New Master is %s ", cluster.master.URL)
	if SaveMasterURL == cluster.master.URL {
		cluster.LogPrintf("ERROR : same server URL after switchover")
		return false
	}
	for _, s := range cluster.slaves {
		cluster.LogPrintf("INFO : Server  %s is %s", s.URL, s.ReadOnly)
		if s.ReadOnly != "OFF" {
			cluster.closeTestCluster(conf, test)
			return false
		}
	}
	cluster.closeTestCluster(conf, test)
	return true
}