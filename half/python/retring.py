'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 
    Usage: 
-------------------------------------------------
'''

from tenacity import retry
import asyncio
import aiohttp


@retry()
async def never_give_up_never_surrender():
    print("Retry forever ignoring Exceptions, don't wait between retries")
    asyncio.sleep(1)
    # async with aiohttp.ClientSession() as se:
    #     async with se.get('http://localhost.com', timeout=0.5) as resp:
    #         await resp.text()
    # raise Exception


loop = asyncio.get_event_loop()
loop.run_until_complete(never_give_up_never_surrender())
