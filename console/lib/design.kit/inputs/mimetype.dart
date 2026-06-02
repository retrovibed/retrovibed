import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';

typedef MimetypePreset = ({Widget label, String value});

class Mimetype extends StatefulWidget {
  final String value;
  final ValueChanged<String> onChanged;
  final List<MimetypePreset> presets;
  final InputDecoration decoration;

  const Mimetype({
    super.key,
    required this.value,
    required this.onChanged,
    this.presets = const [],
    this.decoration = const InputDecoration(),
  });

  @override
  State<Mimetype> createState() => _MimetypeState();
}

class _MimetypeState extends State<Mimetype> {
  bool _expanded = false;
  int _generation = 0;

  void _onTextChanged(String text) {
    widget.onChanged(text);
  }

  void _selectPreset(MimetypePreset preset) {
    setState(() {
      _expanded = false;
      _generation++;
    });
    widget.onChanged(preset.value);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      spacing: defaults.spacing,
      children: [
        TextFormField(
          key: ValueKey(_generation),
          initialValue: widget.value,
          decoration: widget.decoration.copyWith(
            isDense: true,
            suffixIcon:
                widget.presets.isEmpty
                    ? widget.decoration.suffixIcon
                    : IconButton(
                      onPressed: () => setState(() => _expanded = !_expanded),
                      icon: Icon(
                        _expanded ? Icons.arrow_drop_up : Icons.arrow_drop_down,
                      ),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(),
                    ),
          ),
          onChanged: _onTextChanged,
        ),
        if (_expanded && widget.presets.isNotEmpty)
          Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            spacing: defaults.spacing,
            children: [
              Container(
                width: double.infinity,
                decoration: BoxDecoration(
                  border: Border(
                    bottom: BorderSide(color: Theme.of(context).dividerColor),
                  ),
                ),
                child: const Text('presets'),
              ),
              LayoutBuilder(
                builder: (context, constraints) {
                  return Wrap(
                    spacing: defaults.spacing,
                    runSpacing: defaults.spacing,
                    alignment: WrapAlignment.start,
                    children:
                        widget.presets
                            .map(
                              (p) => OutlinedButton(
                                onPressed: () => _selectPreset(p),
                                style: OutlinedButton.styleFrom(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 16,
                                    vertical: 8,
                                  ),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                ),
                                child: p.label,
                              ),
                            )
                            .toList(),
                  );
                },
              ),
            ],
          ),
      ],
    );
  }
}
