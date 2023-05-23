# gpq

A utility for working with [GeoParquet](https://github.com/opengeospatial/geoparquet).

## Installation

The `gpq` program can be installed by downloading one of the archives from [the latest release](https://github.com/planetlabs/gpq/releases).

Extract the archive and place the `gpq` executable somewhere on your path.  See a list of available commands by running `gpq` in your terminal.

For Homebrew users, you can install `gpq` from the [Planet tap](https://github.com/planetlabs/homebrew-tap):

```shell
# run `brew update` first if you have used this tap previously and want the latest formula
brew install planetlabs/tap/gpq
```

## WebAssembly

In addition to the CLI program, the `gpq` utility is built as a WebAssembly binary.  The WASM build can be downloaded from [the latest release](https://github.com/planetlabs/gpq/releases).

To give it a try without downloading or installing anything, see https://planetlabs.github.io/gpq/.

## Command Line Utility

The `gpq` program can be used to validate GeoParquet files and to convert to and from GeoJSON.

```shell
# see the available commands
gpq --help
```

### validate

The `validate` command validates the "geo" file metadata against [the schema](https://github.com/opengeospatial/geoparquet/blob/main/format-specs/schema.json).

```shell
gpq validate example.parquet
```

In the future, this command might also read geometries to confirm that they are valid and that any provided `bbox` is correct.  But for now, just the "geo" metadata is validated against the schema.

### convert

The `convert` command can convert a GeoJSON file to GeoParquet or a GeoParquet file to GeoJSON.

```shell
# read geojson and write geoparquet
gpq convert example.geojson example.parquet
```

```shell
# read geoparquet and write geojson
gpq convert example.parquet example.geojson
```

### describe

The `describe` command prints schema information and metadata about a GeoParquet file.

```shell
# use the `--pretty` argument to format the JSON output
gpq describe example.parquet
```

## Limitations

 * Non-geographic CRS information is not preserved when converting GeoParquet to GeoJSON.
 * Page and row group size is not configurable when writing GeoParquet.  This may change soon.
 * Reading GeoParquet files with multiple geometry columns is supported.  Reading GeoJSON files with multiple geometry properties is not supported.
 * Feature identifiers in GeoJSON are not written to GeoParquet columns.  This may change soon.
