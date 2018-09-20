package dataStore

func DerivedResults(qr Results, ch <-chan Result) Results {
	return &results{
		query:   qr.Query(),
		process: qr.Process(),
		result:  ch,
	}
}

func NaiveFilter(qr Results, filter Filter) Results {
	ch := make(chan Result)
	go func() {
		defer close(ch)
		defer qr.Close()

		for e := range qr.Next() {
			if e.Error != nil || filter.Filter(e.entry) {
				ch <- e
			}
		}
	}()

	return DerivedResults(qr, ch)
}

// NaiveLimit truncates the results to a given int limit
func NaiveLimit(qr Results, limit int) Results {
	ch := make(chan Result)
	go func() {
		defer close(ch)
		defer qr.Close()

		l := 0
		for e := range qr.Next() {
			if e.Error != nil {
				ch <- e
				continue
			}
			ch <- e
			l++
			if limit > 0 && l >= limit {
				break
			}
		}
	}()

	return DerivedResults(qr, ch)
}

func NaiveOffset(qr Results, offset int) Results {
	ch := make(chan Result)
	go func() {
		defer close(ch)
		defer qr.Close()

		sent := 0
		for e := range qr.Next() {
			if e.Error != nil {
				ch <- e
			}

			if sent < offset {
				sent++
				continue
			}
			ch <- e
		}
	}()

	return DerivedResults(qr, ch)
}

func NaiveOrder(qr Results, o Order) Results {
	ch := make(chan Result)
	var entries []Entry
	go func() {
		defer close(ch)
		defer qr.Close()

		for e := range qr.Next() {
			if e.Error != nil {
				ch <- e
			}

			entries = append(entries, e.entry)
		}

		o.Sort(entries)
		for _, e := range entries {
			ch <- Result{entry: e}
		}
	}()

	return DerivedResults(qr, ch)
}

func NaiveQueryApply(q Query, qr Results) Results {
	if q.Prefix != "" {
		qr = NaiveFilter(qr, FilterKeyPrefix{q.Prefix})
	}
	for _, f := range q.Filters {
		qr = NaiveFilter(qr, f)
	}
	for _, o := range q.Orders {
		qr = NaiveOrder(qr, o)
	}
	if q.Offset != 0 {
		qr = NaiveOffset(qr, q.Offset)
	}
	if q.Limit != 0 {
		qr = NaiveLimit(qr, q.Offset)
	}
	return qr
}

func ResultEntriesFrom(keys []string, vals [][]byte) []Entry {
	re := make([]Entry, len(keys))
	for i, k := range keys {
		re[i] = Entry{Key: k, Value: vals[i]}
	}
	return re
}
