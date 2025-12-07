package cli

import (
	"fmt"
	"strings"
	"unicode"
)

type benchSources struct {
	Main string
}

// safeIdent converts a project name to a valid C++ identifier (replaces hyphens with underscores)
func safeIdent(name string) string {
	if name == "" {
		return "project"
	}
	var b strings.Builder
	for i, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			if i == 0 && unicode.IsDigit(r) {
				b.WriteByte('_')
			}
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "project"
	}
	return b.String()
}

func generateBenchmarkArtifacts(projectName string, bench string) (*benchSources, []string) {
	switch bench {
	case "google-benchmark":
		return &benchSources{Main: googleBenchMain(projectName)}, []string{"benchmark"}
	case "nanobench":
		return &benchSources{Main: nanoBenchMain(projectName)}, []string{"nanobench"}
	case "catch2-benchmark":
		return &benchSources{Main: catch2BenchMain(projectName)}, []string{"catch2"}
	default:
		return nil, nil
	}
}

func googleBenchMain(projectName string) string {
	safeName := safeIdent(projectName)
	return fmt.Sprintf(`#include <benchmark/benchmark.h>
#include <%s/%s.hpp>

static void BM_version(benchmark::State& state) {
    for (auto _ : state) {
        benchmark::DoNotOptimize(%s::version());
    }
}

BENCHMARK(BM_version);

int main(int argc, char** argv) {
    benchmark::Initialize(&argc, argv);
    if (benchmark::ReportUnrecognizedArguments(argc, argv)) return 1;
    benchmark::RunSpecifiedBenchmarks();
}
`, projectName, projectName, safeName)
}

func nanoBenchMain(projectName string) string {
	safeName := safeIdent(projectName)
	return fmt.Sprintf(`#include <nanobench.h>
#include <%s/%s.hpp>
#include <iostream>

int main() {
    ankerl::nanobench::Bench bench;
    bench.run("version", [] {
        ankerl::nanobench::doNotOptimizeAway(%s::version());
    });
    return 0;
}
`, projectName, projectName, safeName)
}

func catch2BenchMain(projectName string) string {
	safeName := safeIdent(projectName)
	return fmt.Sprintf(`#include <catch2/catch_all.hpp>
#include <%s/%s.hpp>

TEST_CASE("Benchmark version", "[benchmark]") {
    BENCHMARK("version") {
        return %s::version();
    };
}
`, projectName, projectName, safeName)
}
