package world

type ActorTime struct {
	Time            uint64
	Actor_id        ActorId
	Stability_index uint64 // To ensure stable sorting.
}

type ActorSchedule struct {
	Actor_times          []ActorTime
	Next_stability_index uint64 // To ensure stable sorting.
}

func MakeActorSchedule() ActorSchedule {
	return ActorSchedule{
		Actor_times: make([]ActorTime, 0, 8),
	}
}

// Implement sort.Interface.
func (self ActorSchedule) Len() int {
	return len(self.Actor_times)
}
func (self ActorSchedule) Less(i, j int) bool {
	// Order by time has priority.
	if self.Actor_times[i].Time < self.Actor_times[j].Time {
		return true
	}
	// But in case of tie (two actors have the same time), the first actor
	// scheduled wins.
	return self.Actor_times[i].Stability_index < self.Actor_times[j].Stability_index
}
func (self ActorSchedule) Swap(i, j int) {
	self.Actor_times[i], self.Actor_times[j] = self.Actor_times[i], self.Actor_times[j]
}

func (self ActorSchedule) Copy() ActorSchedule {
	// `self` is already a copy.  We just need to copy the content of the slice.
	new_slice := make([]ActorTime, len(self.Actor_times))
	copy(new_slice, self.Actor_times)
	self.Actor_times = new_slice
	return self
}

// Find an actor that should be acting at the provided time.
// Brute force search that does ot assume that the actors are sorted by time.
func (self ActorSchedule) Next(time uint64) (ActorTime, bool) {
	for _, actor_time := range self.Actor_times {
		if actor_time.Time <= time {
			return actor_time, true
		}
	}
	return ActorTime{}, false
}

func (self ActorSchedule) PosActorId(actor_id ActorId) int {
	for index, actor_time := range self.Actor_times {
		if actor_time.Actor_id == actor_id {
			return index
		}
	}
	return -1
}

func (self ActorSchedule) Pos(actor_time ActorTime) int {
	for index, actor_time0 := range self.Actor_times {
		if actor_time0 == actor_time {
			return index
		}
	}
	return -1
}

func (self ActorSchedule) Remove(actor_time ActorTime) (ActorSchedule, bool) {
	index := self.Pos(actor_time)
	if index == -1 {
		return self, false
	}
	// A cheap remove consists in swapping the last and the current, then reduce
	// the length.  I would avoid that for two reasons.
	// 1: my slice is a priori shared, so I need an expensive  deep copy anyway.
	// 2: I would like to keep the order intact in order to speed up search-by
	//    time.

	// I call Copy in order to deep-copy the Actor_times slice.
	// If I do not perform a deep copy of the slice, then the call to `append`
	// later will affect the content of the original slice, introducing side
	// effects.
	self = self.Copy()
	self.Actor_times = append(
		self.Actor_times[:index],
		self.Actor_times[index+1:]...)
	return self, true
}

func (self ActorSchedule) Add(actor_id ActorId, time uint64) ActorSchedule {
	new_entry := ActorTime{
		Actor_id:        actor_id,
		Time:            time,
		Stability_index: self.Next_stability_index,
	}
	len_slice := len(self.Actor_times)
	result := ActorSchedule{
		Next_stability_index: self.Next_stability_index + 1,
		Actor_times:          make([]ActorTime, len_slice, len_slice+1),
	}
	copy(result.Actor_times, self.Actor_times)
	result.Actor_times = append(result.Actor_times, new_entry)
	return result
}
