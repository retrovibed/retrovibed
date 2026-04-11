import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;

enum VideoResolution {
  Any('Any', 0),
  r480p('480p', 480),
  r720p('720p', 720),
  r1080p('1080p', 1080),
  r1440p('1440p (QHD)', 1440),
  r2160p('2160p (4K)', 2160);

  final String label;
  final int height;

  const VideoResolution(this.label, this.height);
}

class VideoSettings extends StatefulWidget {
  const VideoSettings({super.key, this.margin = EdgeInsets.zero});
  final EdgeInsets margin;

  @override
  State<VideoSettings> createState() => _VideoSettingsState();
}

class _VideoSettingsState extends State<VideoSettings> {
  VideoResolution _resolution = VideoResolution.Any;

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    List<DropdownMenuItem<VideoResolution>> buildDropdownItems() {
      return VideoResolution.values.map((resolution) {
        return DropdownMenuItem<VideoResolution>(
          value: resolution,
          child: Text(resolution.label),
        );
      }).toList();
    }

    return ds.Container(
      padding: defaults.padding,
      margin: widget.margin,
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing / 2.5,
        children: [
          forms.Field(
            label: Text("resolution"),
            input: DropdownButton(
              alignment: Alignment.topLeft,
              isExpanded: true,
              value: _resolution,
              items: buildDropdownItems(),
              onChanged: (VideoResolution? v) {
                if (v == null) return;
                setState(() {
                  _resolution = v;
                });
              },
            ),
          ),
        ],
      ),
    );
  }
}
