import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import './api.dart' as api;
import './icons.dart' as wgicons;

class Edit extends StatelessWidget {
  final api.Wireguard current;
  final Function(api.Wireguard)? onChange;

  Edit(this.current, {super.key, this.onChange});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          forms.Field(
            label: Text("description"),
            input: Row(
              children: [
                Expanded(
                  child: TextFormField(
                    initialValue: current.description,
                    maxLines: 1,
                    onChanged: (v) => onChange?.call(current..description = v),
                  ),
                ),
                wgicons.Icons.delete(
                  current,
                  onPressed: () async {
                    final modal = ds.modals.of(context);
                    modal?.push(
                      ds.Confirmation.yesNo(
                        content: Text('Delete ${current.description}?'),
                        onCancel: (_) => modal.push(null),
                        onConfirm: (_) {
                          api.wireguard.delete(current.id).then((_) {
                            modal.push(null);
                          });
                        },
                      ),
                    );
                  },
                ),
              ],
            ),
          ),
          forms.Field(
            label: Text("port"),
            input: TextFormField(
              decoration: InputDecoration(
                helperText: "0 means automatic port detection",
              ),
              initialValue: current.port.toString(),
              keyboardType: TextInputType.number,
              onChanged: (v) {
                final port = int.tryParse(v);
                if (port != null) {
                  onChange?.call(current..port = port);
                }
              },
            ),
          ),
          forms.Field(
            label: Text("dns rate limit"),
            input: inputs.RateLimit(
              value: current.dnsRateLimit,
              presets: const [
                (label: '10/sec', value: 10, unit: 'sec'),
                (label: '100/sec', value: 100, unit: 'sec'),
                (label: '1000/sec', value: 1000, unit: 'sec'),
              ],
              onChanged: (v) {
                onChange?.call(current..dnsRateLimit = v);
              },
            ),
          ),
          forms.Field(
            label: Text("maximum connections"),
            input: inputs.Uint64(
              value: current.maximumConnections,
              presets: [
                (label: '16', value: ds.Int64(16)),
                (label: '32', value: ds.Int64(32)),
                (label: '64', value: ds.Int64(64)),
                (label: 'unlimited', value: ds.Int64(-1)),
              ],
              onChanged: (v) {
                onChange?.call(current..maximumConnections = v);
              },
            ),
          ),
        ],
      ),
    );
  }
}
