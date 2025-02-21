import 'package:flutter/material.dart';
import 'package:fractal/design.kit/forms.dart' as forms;

class SettingsLeech extends StatefulWidget {
  SettingsLeech({super.key});

  @override
  State<SettingsLeech> createState() => _EditView();
}

class _EditView extends State<SettingsLeech> {
  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
      constraints: BoxConstraints(minHeight: 32, minWidth: 32),
      child: forms.Container(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            forms.Field(
              label: Text("ratio"),
              input: TextFormField(
                decoration: new InputDecoration(
                  hintText: "2.0",
                  helperText: "seed until this ratio, default is unlimited",
                ),
                keyboardType: TextInputType.number,
                // onChanged: (v) => widget.feed.description = v,
              ),
            ),
            forms.Field(
              label: Text("upload"),
              input: TextFormField(
                decoration: new InputDecoration(
                  hintText: "0",
                  helperText:
                      "maximum upload rate allowed per torrent (KB/s), default is unlimited",
                ),
                keyboardType: TextInputType.number,
                // onChanged: (v) => widget.feed.description = v,
              ),
            ),
            forms.Field(
              label: Text("download"),
              input: TextFormField(
                decoration: new InputDecoration(
                  hintText: "0",
                  helperText:
                      "maximum upload rate allowed per torrent (KB/s), default is unlimited",
                ),
                keyboardType: TextInputType.number,
                // onChanged: (v) => widget.feed.description = v,
              ),
            ),
            forms.Field(
              label: Text("peers"),
              input: TextFormField(
                decoration: new InputDecoration(
                  hintText: "128",
                  helperText:
                      "maximum number of peers allowed per torrent, default is 128",
                ),
                keyboardType: TextInputType.number,
                // onChanged: (v) => widget.feed.description = v,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
