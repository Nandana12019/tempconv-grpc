import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

void main() {
  runApp(const TempConvApp());
}

class TempConvApp extends StatelessWidget {
  const TempConvApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'TempConv',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepOrange),
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  // REST Gateway (GKE LoadBalancer)
  String get _apiBase {
    return 'http://35.223.167.112:8080/api';
  }

  String? _error;
  String? _result;

  Future<void> _celsiusToFahrenheit(String celsiusStr) async {
    _error = null;
    _result = null;
    setState(() {});

    final c = double.tryParse(celsiusStr);
    if (c == null) {
      setState(() => _error = 'Enter a valid number for Celsius');
      return;
    }

    try {
      final uri = Uri.parse('$_apiBase/c2f');

      final response = await http.post(
        uri,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'celsius': c}),
      );

      final r = jsonDecode(response.body) as Map<String, dynamic>;
      final f = (r['fahrenheit'] as num).toDouble();

      setState(() =>
          _result = '${c.toStringAsFixed(1)} °C = ${f.toStringAsFixed(1)} °F');
    } catch (e) {
      setState(() => _error = 'Network error: $e');
    }
  }

  Future<void> _fahrenheitToCelsius(String fahrenheitStr) async {
    _error = null;
    _result = null;
    setState(() {});

    final f = double.tryParse(fahrenheitStr);
    if (f == null) {
      setState(() => _error = 'Enter a valid number for Fahrenheit');
      return;
    }

    try {
      final uri = Uri.parse('$_apiBase/f2c');

      final response = await http.post(
        uri,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'fahrenheit': f}),
      );

      final r = jsonDecode(response.body) as Map<String, dynamic>;
      final c = (r['celsius'] as num).toDouble();

      setState(() =>
          _result = '${f.toStringAsFixed(1)} °F = ${c.toStringAsFixed(1)} °C');
    } catch (e) {
      setState(() => _error = 'Network error: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    final cController = TextEditingController();
    final fController = TextEditingController();

    return Scaffold(
      appBar: AppBar(
        title: const Text('TempConv'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const Text(
                  'Celsius → Fahrenheit',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 8),
                TextField(
                  controller: cController,
                  keyboardType:
                      const TextInputType.numberWithOptions(decimal: true),
                  decoration: const InputDecoration(
                    hintText: 'e.g. 100',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 8),
                FilledButton(
                  onPressed: () =>
                      _celsiusToFahrenheit(cController.text.trim()),
                  child: const Text('Convert to °F'),
                ),
                const SizedBox(height: 32),
                const Text(
                  'Fahrenheit → Celsius',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 8),
                TextField(
                  controller: fController,
                  keyboardType:
                      const TextInputType.numberWithOptions(decimal: true),
                  decoration: const InputDecoration(
                    hintText: 'e.g. 212',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 8),
                FilledButton(
                  onPressed: () =>
                      _fahrenheitToCelsius(fController.text.trim()),
                  child: const Text('Convert to °C'),
                ),
                if (_error != null) ...[
                  const SizedBox(height: 24),
                  Text(_error!,
                      style: TextStyle(
                          color: Theme.of(context).colorScheme.error)),
                ],
                if (_result != null) ...[
                  const SizedBox(height: 24),
                  Text(_result!,
                      style: Theme.of(context).textTheme.titleMedium),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }
}

