package typeflow

import (
    "bufio"
    "io"
    "os"
    "strings"
    "testing"
)

type word_source_test struct {
    substr           string
    expected_matches []expected_match
}

var word_source_tests = []word_source_test{
    {"rep of ireland", []expected_match{{"ireland (republic)", similarity_range{0.32, 0.34}}}},
    {"united", []expected_match{{"united states", similarity_range{0.4, 0.5}}}},
}

func TestWordSource(t *testing.T) {
    ws := NewWordSource()

    // building country name
    // source from file
    file, err := os.Open("testdata/countries.txt")
    country_names := make([]string, 0)
    if err != nil {
        t.Log("Cannot open expected file testdata/countries.txt. Skipping this test.")
        t.SkipNow()
        return
    }
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        }
        country_names = append(country_names, strings.ToLower(line[:len(line)-1]))
    }

    ws.SetSource(country_names)

OuterLoop:
    for _, test := range word_source_tests {
        t.Logf("Finding matches for substring '%s'", test.substr)
        matches, err := ws.FindMatch(test.substr, 0.32)
        if err != nil {
            t.Logf("An error occurred: %v", err)
            for _, m := range matches {
                t.Logf("%s, %f", m.String, m.Similarity)
            }
            t.FailNow()
        }

        for _, match := range matches {
            for _, expected := range test.expected_matches {
                if match.String == expected.string && match.Similarity >= expected.similarity_range.low &&
                    match.Similarity <= expected.similarity_range.high {
                    t.Log("Found!")
                    continue OuterLoop
                }
            }
        }
        t.Logf("Couldn't find expected match")
        t.Logf("Found the following matches:")
        for _, m := range matches {
            t.Logf("'%s', '%f'", m.String, m.Similarity)
        }
        t.FailNow()
    }
}

func BenchmarkFindMatch(b *testing.B) {
    b.StopTimer()
    ws := NewWordSource()
    // building country name
    // source from file
    file, err := os.Open("testdata/countries.txt")
    country_names := make([]string, 0)
    if err != nil {
        b.Log("Cannot open expected file testdata/countries.txt. Skipping this test.")
        b.SkipNow()
        return
    }
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        }
        country_names = append(country_names, strings.ToLower(line[:len(line)-1]))
    }

    ws.SetSource(country_names)

    b.ResetTimer()
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        test := word_source_tests[i%len(word_source_tests)]
        matches, err := ws.FindMatch(test.substr, 0.32)
        if err != nil {
            b.Errorf("An error occurred: %v", err)
            for _, m := range matches {
                b.Logf("%s, %f", m.String, m.Similarity)
            }
        }
    }
}
