package cli

import (
	"fmt"
	"strings"
)

type benchSources struct {
	Main string
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
	return fmt.Sprintf(`#include <benchmark/benchmark.h>
#include <%s/%s.hpp>

static void BM_greet(benchmark::State& state) {
    for (auto _ : state) {
        benchmark::DoNotOptimize(%s::version());
    }
}

BENCHMARK(BM_greet);

int main(int argc, char** argv) {
    benchmark::Initialize(&argc, argv);
    if (benchmark::ReportUnrecognizedArguments(argc, argv)) return 1;
    benchmark::RunSpecifiedBenchmarks();
}
`, projectName, projectName, projectName)
}

func nanoBenchMain(projectName string) string {
	return fmt.Sprintf(`#include <nanobench.h>
#include <%s/%s.hpp>
#include <iostream>

int main() {
    ankerl::nanobench::Bench bench;
    bench.run("greet", [] {
        ankerl::nanobench::doNotOptimizeAway(%s::version());
    });
    return 0;
}
`, projectName, projectName, projectName)
}

func catch2BenchMain(projectName string) string {
	return fmt.Sprintf(`#define CATCH_CONFIG_MAIN
#include <catch2/catch_all.hpp>
#include <%s/%s.hpp>

TEST_CASE("benchmark greet", "[bench]") {
    auto& bench = Catch::getResultCapture().benchmark("greet");
    bench.measure([] { return %s::version(); });
}
`, projectName, projectName, projectName)
}

func containsIgnoreCase(haystack []string, needle string) bool {
	needle = strings.ToLower(needle)
	for _, s := range haystack {
		if strings.ToLower(s) == needle {
			return true
		}
	}
	return false
}
