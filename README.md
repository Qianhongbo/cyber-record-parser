# cyber-record-parser

A Go library for parsing and processing Cyber Record files.

## Build

```bash
go build -o cyber_record_parser ./cmd
```

Move the binary to a location in your PATH.

```bash
mv cyber_record_parser /usr/local/bin
```

## Usage

```bash
# help
cyber_record_parser -h

# info cmd
cyber_record_parser info <file>

# echo cmd
cyber_record_parser echo <file> -t/--topic <topic>
```

> Note:
>
> Pause and resume echo cmd with `Space` key.
>
> Exit echo cmd with `Ctrl+C` / `ESC` / `q`.

## Output example

### Info cmd

```bash
cyber_record_parser info /path/to/file
```

```Text

Cyber Record information:
----------------------------

- Record file path     sensor_rgb.record
- Version              1.0
- Size                 7.6 GB
- Compression          COMPRESS_NONE
- Chunk raw size       210 MB
- Chunk interval       20s
- Start time           2018-01-03 19:37:30 +0800 CST
- End time             2018-01-03 19:38:31 +0800 CST
- Duration             1m1s
- Message number       43820
- Channel number       14
- Is complete          true

Channels information:
----------------------------

Channel name                                       | Count   | Type
/apollo/canbus/chassis                             | 5851    | apollo.canbus.Chassis
/apollo/localization/pose                          | 5856    | apollo.localization.LocalizationEstimate
/apollo/sensor/camera/traffic/image_long           | 471     | apollo.drivers.Image
/apollo/sensor/camera/traffic/image_short          | 469     | apollo.drivers.Image
/apollo/sensor/conti_radar                         | 789     | apollo.drivers.ContiRadar
/apollo/sensor/gnss/best_pose                      | 59      | apollo.drivers.gnss.GnssBestPose
/apollo/sensor/gnss/imu                            | 11630   | apollo.drivers.gnss.Imu
/apollo/sensor/gnss/ins_stat                       | 118     | apollo.drivers.gnss.InsStat
/apollo/sensor/gnss/odometry                       | 5848    | apollo.localization.Gps
/apollo/sensor/gnss/rtk_eph                        | 49      | apollo.drivers.gnss.GnssEphemeris
/apollo/sensor/gnss/rtk_obs                        | 352     | apollo.drivers.gnss.EpochObservation
/apollo/sensor/velodyne64/compensator/PointCloud2  | 587     | apollo.drivers.PointCloud
/tf                                                | 11740   | apollo.transform.TransformStampeds
/tf_static                                         | 1       | apollo.transform.TransformStampeds
```

### Echo cmd

```bash
cyber_record_parser echo /path/to/file -t <topic>
```

```Text
--------------------------------------------------
Channel name: /asensing_rtk570/raw_rtk_can
Time nanosecond: 1733119089389686437
Time: 2024-12-02 13:58:09

Message:
{

}
```

### tojson cmd

```bash
cyber_record_parser tojson /path/to/file -t <topic> -o <output_file>
```

```Text
Save topic (/asensing_rtk570/raw_rtk_can) messages to /tmp/test.json
```

## Use as a library

```go
package main

import (
	"fmt"

	"cyber_record_parser/internal/record"
)

func main() {
	r := record.NewRecord("/path/to/file")
	defer r.Close()

	for msg := range r.ReadMessages() {
		fmt.Printf("Channel: %s\n", msg.ChannelName)
		fmt.Printf("Timestamp: %d\n", msg.NanoTimestamp)
		fmt.Printf("Content: %s\n", string(msg.Content))
	}
}
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
