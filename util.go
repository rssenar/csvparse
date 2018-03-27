package csvparse

// for _, row := range rows[:][1:] {
// 	w := csv.NewWriter(os.Stdout)
// 	if err := w.Write(row); err != nil {
// 		return fmt.Errorf("%v : error writing record to csv", err)
// 	}
// 	// Write any buffered data to the underlying writer (standard output).
// 	w.Flush()

// 	if err := w.Error(); err != nil {
// 		return fmt.Errorf("%v : error writing to stdout", err)
// 	}
// }

// for a, b := range rows[:][0] {
// 	fmt.Println(a, b)
// }

// type kv struct {
// 	Key   string
// 	Value int
// }

// var ss []kv
// for k, v := range p.header {
// 	ss = append(ss, kv{k, v})
// }

// sort.Slice(ss, func(i, j int) bool {
// 	return ss[i].Value < ss[j].Value
// })

// for _, str := range ss {
// 	fmt.Printf("%s, %d\n", str.Key, str.Value)
// }
// fmt.Printf("%v", x)
