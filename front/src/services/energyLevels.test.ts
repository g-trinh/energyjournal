import { getEnergyLevels, saveEnergyLevels, type EnergyLevels } from './energyLevels'

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('energyLevels service', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('returns parsed levels for successful GET', async () => {
    const payload: EnergyLevels = {
      date: '2026-02-21',
      physical: 7,
      mental: 5,
      emotional: 8,
    }
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(200, payload),
    )
    const abortController = new AbortController()

    const result = await getEnergyLevels(
      '2026-02-21',
      'token',
      abortController.signal,
    )

    expect(result).toEqual(payload)
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/energy/levels?date=2026-02-21',
      expect.objectContaining({
        method: 'GET',
        signal: abortController.signal,
      }),
    )
  })

  it('returns null for 404 GET', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(404, { error: 'not found' }),
    )

    const result = await getEnergyLevels('2026-02-21', 'token')
    expect(result).toBeNull()
  })

  it('throws for non-404 GET failures', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(500, { error: 'boom' }),
    )

    await expect(getEnergyLevels('2026-02-21', 'token')).rejects.toThrow()
  })

  it('throws on GET network failure', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('network down'))

    await expect(getEnergyLevels('2026-02-21', 'token')).rejects.toThrow(
      'network down',
    )
  })

  it('returns parsed levels for successful PUT', async () => {
    const payload: EnergyLevels = {
      date: '2026-02-21',
      physical: 6,
      mental: 5,
      emotional: 4,
    }
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(jsonResponse(200, payload))

    const result = await saveEnergyLevels(payload, 'token')
    expect(result).toEqual(payload)
  })

  it('throws for non-200 PUT failures', async () => {
    const payload: EnergyLevels = {
      date: '2026-02-21',
      physical: 6,
      mental: 5,
      emotional: 4,
    }
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(400, { error: 'invalid' }),
    )

    await expect(saveEnergyLevels(payload, 'token')).rejects.toThrow()
  })
})
