import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

class DownloadRowDisplay extends StatelessWidget {
  final api.Download current;
  final Widget? Function(BuildContext)? trailing;
  final Widget help;
  const DownloadRowDisplay({super.key, required this.current, this.trailing, this.help = ds.HelpScope.None});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    var percentage = math.min(
      current.bytes == 0 ? 0.0 : (current.downloaded.toDouble() / current.bytes.toDouble()),
      1.0,
    );

    return ds.Help(
      ds.ErrorBoundary(
        SelectionArea(
          child: Builder(
            builder: (context) {
              final compact = defaults.isCompact;
              if (compact) {
                return ds.TableRow.single(
                  padding: defaults.padding,
                  Column(
                    children: [
                      Row(
                        spacing: defaults.spacing,
                        children: [
                          Expanded(
                            child: Text(
                              current.media.description,
                              overflow: TextOverflow.ellipsis,
                              maxLines: 1,
                            ),
                          ),
                          Expanded(
                            child: LinearProgressIndicator(
                              value: percentage,
                              semanticsLabel: 'Linear progress indicator',
                            ),
                          ),
                        ],
                      ),
                      Row(
                        spacing: defaults.spacing,
                        children: [
                          Expanded(child: Icon(Icons.people_outline, size: 16)),
                          Expanded(
                            child: Text(
                              current.peers.toString().padLeft(3),
                              style: const TextStyle(fontFamily: 'monospace'),
                            ),
                          ),
                          Expanded(
                            child: Text(
                              "${(percentage * 100).toStringAsFixed(2).padLeft(6)}%",
                              style: const TextStyle(fontFamily: 'monospace'),
                            ),
                          ),
                          trailing?.call(context) ?? const SizedBox(),
                        ],
                      ),
                    ],
                  ),
                );
              }

              return ds.TableRow(padding: defaults.padding, [
                Expanded(
                  child: Text(
                    current.media.description,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                Expanded(
                  child: LinearProgressIndicator(
                    value: percentage,
                    semanticsLabel: 'Linear progress indicator',
                  ),
                ),
                Icon(Icons.people_outline, size: 16),
                Text(current.peers.toString().padLeft(3), style: const TextStyle(fontFamily: 'monospace')),
                Text(
                  "${(percentage * 100).toStringAsFixed(2).padLeft(6)}%",
                  style: const TextStyle(fontFamily: 'monospace'),
                ),
                trailing?.call(context) ?? const SizedBox(),
              ]);
            },
          ),
        ),
      ),
      help,
    );
  }
}
