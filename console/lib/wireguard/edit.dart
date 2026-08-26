import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;
import './icons.dart' as wgicons;

class Edit extends StatefulWidget {
  final api.Wireguard current;
  final Function(api.Wireguard)? onChange;

  const Edit(this.current, {super.key, this.onChange});

  @override
  State<Edit> createState() => _EditState();
}

class _EditState extends State<Edit> with LoadingState {
  late api.Wireguard _current = widget.current;

  @override
  void initState() {
    super.initState();
    loading = false;
  }

  @override
  void didUpdateWidget(covariant Edit oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (!identical(oldWidget.current, widget.current)) {
      _current = widget.current;
    }
  }

  void _onChange(api.Wireguard updated) {
    setState(() {
      _current = updated;
      loading = true;
    });

    Future.sync(() => widget.onChange?.call(updated))
        .then((_) {
          setState(() => loading = false);
        })
        .catchError((error) {
          setState(() {
            loading = false;
            cause = ds.Errors.httpauto(error, onTap: resetCause);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((error) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(error, onTap: resetCause);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final current = _current;
    return forms.Container(
      cause: cause,
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
                    onChanged: (v) => _onChange(current..description = v),
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
                  _onChange(current..port = port);
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
                _onChange(current..dnsRateLimit = v);
              },
            ),
          ),
          forms.Field(
            label: Text("outbound rate limit"),
            input: inputs.RateLimit(
              value: current.outboundRateLimit,
              presets: const [
                (label: '4/sec', value: 4, unit: 'sec'),
                (label: '10/sec', value: 10, unit: 'sec'),
                (label: '100/sec', value: 100, unit: 'sec'),
              ],
              onChanged: (v) {
                _onChange(current..outboundRateLimit = v);
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
                _onChange(current..maximumConnections = v);
              },
            ),
          ),
          const Divider(),
          forms.Presets<api.Wireguard>(
            current: current,
            presets: [
              (label: 'ProtonVPN', apply: (w) => w..outboundRateLimit = 2),
            ],
            onSelected: _onChange,
          ),
        ],
      ),
    );
  }
}

