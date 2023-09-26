package command_test

import (
	"encoding/json"

	"github.com/planetlabs/gpq/cmd/gpq/command"
	"github.com/planetlabs/gpq/internal/test"
)

func (s *Suite) TestDescribe() {
	cmd := &command.DescribeCmd{
		Input:  "../../../internal/testdata/cases/example-v1.0.0.parquet",
		Format: "json",
	}

	s.Require().NoError(cmd.Run())

	output := s.readStdout()
	info := &command.DescribeInfo{}
	err := json.Unmarshal(output, info)
	s.Require().NoError(err)

	s.Equal(int64(5), info.NumRows)
	s.Require().Len(info.Schema.Fields, 6)

	s.Equal("geometry", info.Schema.Fields[0].Name)
	s.Equal("binary", info.Schema.Fields[0].Type)
	s.Equal("gzip", info.Schema.Fields[0].Compression)
	s.True(info.Schema.Fields[0].Optional)

	s.Equal("pop_est", info.Schema.Fields[1].Name)
	s.Equal("double", info.Schema.Fields[1].Type)
	s.Equal("gzip", info.Schema.Fields[1].Compression)
	s.True(info.Schema.Fields[1].Optional)

	s.Equal("continent", info.Schema.Fields[2].Name)
	s.Equal("binary", info.Schema.Fields[2].Type)
	s.Equal("string", info.Schema.Fields[2].Annotation)
	s.Equal("gzip", info.Schema.Fields[2].Compression)
	s.True(info.Schema.Fields[2].Optional)

	s.Equal("gdp_md_est", info.Schema.Fields[3].Name)
	s.Equal("double", info.Schema.Fields[3].Type)
	s.Equal("gzip", info.Schema.Fields[3].Compression)
	s.True(info.Schema.Fields[3].Optional)

	s.Equal("iso_a3", info.Schema.Fields[4].Name)
	s.Equal("binary", info.Schema.Fields[4].Type)
	s.Equal("string", info.Schema.Fields[4].Annotation)
	s.Equal("gzip", info.Schema.Fields[4].Compression)
	s.True(info.Schema.Fields[4].Optional)

	s.Equal("name", info.Schema.Fields[5].Name)
	s.Equal("binary", info.Schema.Fields[5].Type)
	s.Equal("string", info.Schema.Fields[5].Annotation)
	s.Equal("gzip", info.Schema.Fields[5].Compression)
	s.True(info.Schema.Fields[5].Optional)
}

func (s *Suite) TestDescribeFromStdin() {
	s.writeStdin(test.GeoParquetFromJSON(s.T(), `{
		"type": "FeatureCollection",
		"features": [
			{
				"type": "Feature",
				"properties": {
					"name": "Null Island"
				},
				"geometry": {
					"type": "Point",
					"coordinates": [0, 0]
				}
			}
		]
	}`))

	cmd := &command.DescribeCmd{
		Format: "json",
	}

	s.Require().NoError(cmd.Run())

	output := s.readStdout()
	info := &command.DescribeInfo{}
	err := json.Unmarshal(output, info)
	s.Require().NoError(err)

	s.Equal(int64(1), info.NumRows)
	s.Require().Len(info.Schema.Fields, 2)

	s.Equal("geometry", info.Schema.Fields[0].Name)
	s.Equal("binary", info.Schema.Fields[0].Type)
	s.Equal("zstd", info.Schema.Fields[0].Compression)
	s.True(info.Schema.Fields[0].Optional)

	s.Equal("name", info.Schema.Fields[1].Name)
	s.Equal("binary", info.Schema.Fields[1].Type)
	s.Equal("string", info.Schema.Fields[1].Annotation)
	s.Equal("zstd", info.Schema.Fields[1].Compression)
	s.True(info.Schema.Fields[1].Optional)
}
